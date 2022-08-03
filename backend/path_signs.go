package backend

import (
	"context"
	"encoding/hex"
	"sync"

	"github.com/bloxapp/key-vault/keymanager/models"
	"github.com/ethereum/go-ethereum/common"

	"github.com/prysmaticlabs/prysm/consensus-types/wrapper"

	vault "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/signer"
	slashingprotection "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/backend/store"
)

// Endpoints patterns
const (
	// SignAttestationPattern is the path pattern for sign attestation endpoint
	SignPattern = "accounts/sign"
)

func signsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         SignPattern,
			HelpSynopsis:    "Sign",
			HelpDescription: `Sign`,
			Fields: map[string]*framework.FieldSchema{
				"sign_req": {
					Type:        framework.TypeString,
					Description: "SSZ Serialized sign request object",
					Default:     "",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathSign,
				},
			},
		},
	}
}

func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Load config
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	// Parse request data
	reqEncoded := data.Get("sign_req").(string)
	reqByts, err := hex.DecodeString(reqEncoded)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode sign request hex")
	}

	signReq := &models.SignRequest{}
	if err := b.encoder.Decode(reqByts, signReq); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal sign request")
	}

	var sig []byte
	err = b.lock(signReq.GetPublicKey(), func() error {
		// bring up KeyVault and wallet
		storage := store.NewHashicorpVaultStore(ctx, req.Storage, config.Network)
		options := vault.KeyVaultOptions{}
		options.SetStorage(storage)

		// Open wallet
		kv, err := vault.OpenKeyVault(&options)
		if err != nil {
			return errors.Wrap(err, "failed to open key vault")
		}

		wallet, err := kv.Wallet()
		if err != nil {
			return errors.Wrap(err, "failed to retrieve wallet")
		}

		var (
			protector                           = slashingprotection.NewNormalProtection(storage)
			simpleSigner signer.ValidatorSigner = signer.NewSimpleSigner(wallet, protector, storage.Network())
			sigErr       error
		)

		switch t := signReq.GetObject().(type) {
		case *models.SignRequestBlock:
			wrappedBlk, wrapErr := wrapper.WrappedBeaconBlock(t.Block)
			if err != nil {
				return errors.Wrap(wrapErr, "failed to wrap Phase0 block")
			}
			sig, sigErr = simpleSigner.SignBeaconBlock(wrappedBlk, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestBlockV2:
			wrappedBlk, wrapErr := wrapper.WrappedBeaconBlock(t.BlockV2)
			if err != nil {
				return errors.Wrap(wrapErr, "failed to wrap Altair block")
			}
			sig, sigErr = simpleSigner.SignBeaconBlock(wrappedBlk, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestBlockV3:
			validateErr := validateRequestedFeeRecipient(signReq.PublicKey, config.FeeRecipients, t.BlockV3.Body.ExecutionPayload.FeeRecipient)
			if validateErr != nil {
				return errors.Wrap(validateErr, "refused to sign")
			}
			wrappedBlk, wrapErr := wrapper.WrappedBeaconBlock(t.BlockV3)
			if err != nil {
				return errors.Wrap(wrapErr, "failed to wrap Bellatrix block")
			}
			sig, sigErr = simpleSigner.SignBeaconBlock(wrappedBlk, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestBlindedBlockV3:
			wrappedBlk, wrapErr := wrapper.WrappedBeaconBlock(t.BlindedBlockV3)
			if err != nil {
				return errors.Wrap(wrapErr, "failed to wrap BlindedBellatrix block")
			}
			sig, sigErr = simpleSigner.SignBeaconBlock(wrappedBlk, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestAttestationData:
			sig, sigErr = simpleSigner.SignBeaconAttestation(t.AttestationData, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestSlot:
			sig, sigErr = simpleSigner.SignSlot(t.Slot, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestEpoch:
			sig, sigErr = simpleSigner.SignEpoch(t.Epoch, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestAggregateAttestationAndProof:
			sig, sigErr = simpleSigner.SignAggregateAndProof(t.AggregateAttestationAndProof, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestSyncCommitteeMessage:
			sig, sigErr = simpleSigner.SignSyncCommittee(t.Root, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestSyncAggregatorSelectionData:
			sig, sigErr = simpleSigner.SignSyncCommitteeSelectionData(t.SyncAggregatorSelectionData, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestContributionAndProof:
			sig, sigErr = simpleSigner.SignSyncCommitteeContributionAndProof(t.ContributionAndProof, signReq.SignatureDomain, signReq.PublicKey)
		case *models.SignRequestRegistration:
			validateErr := validateRequestedFeeRecipient(signReq.PublicKey, config.FeeRecipients, t.Registration.FeeRecipient)
			if validateErr != nil {
				return errors.Wrap(validateErr, "refused to sign")
			}
			sig, sigErr = simpleSigner.SignRegistration(t.Registration, signReq.SignatureDomain, signReq.PublicKey)
		default:
			return errors.Errorf("sign request: not supported")
		}

		// Some tests rely on the error message returned by SignBeaconBlock,
		// so this error should not be wrapped!
		return sigErr
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(sig),
		},
	}, nil
}

func (b *backend) lock(pubKeyBytes []byte, cb func() error) error {
	b.signMapLock.Lock()
	pubKey := hex.EncodeToString(pubKeyBytes)
	if _, ok := b.signLock[pubKey]; !ok {
		b.signLock[pubKey] = &sync.Mutex{}
	}
	lock := b.signLock[pubKey]
	b.signMapLock.Unlock()

	lock.Lock()
	err := cb()
	lock.Unlock()

	return err
}

var (
	ErrFeeRecipientNotSet  = errors.New("fee recipient is not configured for public key")
	ErrFeeRecipientDiffers = errors.New("requested fee recipient does not match configured fee recipient")
)

func validateRequestedFeeRecipient(pubKey []byte, configFeeRecipients FeeRecipients, requestedFeeRecipient []byte) error {
	feeRecipient, ok := configFeeRecipients.Get(pubKey)
	if !ok {
		return ErrFeeRecipientNotSet
	}
	if feeRecipient != common.BytesToAddress(requestedFeeRecipient) {
		return ErrFeeRecipientDiffers
	}
	return nil
}
