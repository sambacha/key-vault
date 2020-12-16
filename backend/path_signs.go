package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/signer"
	"github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	v2 "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

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
				"sign_req": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "SSZ Serialized sign request object",
					Default:     "",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSign,
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

	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage, config.Network)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)
	// Open wallet
	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	// Parse request data
	reqEncoded := data.Get("sign_req").(string)
	reqByts, err := hex.DecodeString(reqEncoded)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode sign request hex")
	}

	signReq := &v2.SignRequest{}
	if err := signReq.Unmarshal(reqByts); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal sign request SSZ")
	}

	protector := slashing_protection.NewNormalProtection(storage)
	var simpleSigner signer.ValidatorSigner = signer.NewSimpleSigner(wallet, protector, storage.Network())

	var sig []byte
	switch t := signReq.GetObject().(type) {
	case *v2.SignRequest_Block:
		sig, err = simpleSigner.SignBeaconBlock(t.Block, signReq.SignatureDomain, signReq.PublicKey)
	case *v2.SignRequest_AttestationData:
		sig, err = simpleSigner.SignBeaconAttestation(t.AttestationData, signReq.SignatureDomain, signReq.PublicKey)
	case *v2.SignRequest_Slot:
		sig, err = simpleSigner.SignSlot(t.Slot, signReq.SignatureDomain, signReq.PublicKey)
	case *v2.SignRequest_Epoch:
		sig, err = simpleSigner.SignEpoch(t.Epoch, signReq.SignatureDomain, signReq.PublicKey)
	case *v2.SignRequest_AggregateAttestationAndProof:
		sig, err = simpleSigner.SignAggregateAndProof(t.AggregateAttestationAndProof, signReq.SignatureDomain, signReq.PublicKey)
	default:
		return nil, errors.Errorf("sign request: not supported")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(sig),
		},
	}, nil
}
