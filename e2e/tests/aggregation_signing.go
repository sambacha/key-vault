package tests

import (
	"encoding/hex"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

	"github.com/bloxapp/eth2-key-manager/signer"

	"github.com/bloxapp/eth2-key-manager/core"
	slashingprotection "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AggregationSigning tests sign aggregation endpoint.
type AggregationSigning struct {
}

// Name returns the name of the test.
func (test *AggregationSigning) Name() string {
	return "Test aggregation signing"
}

// Run run the test.
func (test *AggregationSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	// Get wallet
	wallet, err := storage.OpenWallet()
	require.NoError(t, err)

	agg := &ethpb.AggregateAttestationAndProof{
		AggregatorIndex: 0,
		Aggregate: &ethpb.Attestation{
			Data: &ethpb.AttestationData{
				BeaconBlockRoot: make([]byte, 32),
				Target:          &ethpb.Checkpoint{Root: make([]byte, 32)},
				Source:          &ethpb.Checkpoint{Root: make([]byte, 32)},
			},
			Signature:       make([]byte, 96),
			AggregationBits: make([]byte, 1),
		},
		SelectionProof: make([]byte, 96),
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKeyBytes, nil, domain, agg)

	// Sign data
	protector := slashingprotection.NewNormalProtection(inmemory.NewInMemStore(core.PraterNetwork))
	var signer signer.ValidatorSigner = signer.NewSimpleSigner(wallet, protector, storage.Network())

	res, err := signer.SignAggregateAndProof(agg, domain, pubKeyBytes)
	require.NoError(t, err)

	// Send sign attestation request
	sig, err := setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)

	require.EqualValues(t, res, sig)
}

func (test *AggregationSigning) serializedReq(pk, root, domain []byte, agg *ethpb.AggregateAttestationAndProof) (map[string]interface{}, error) {
	req := &validatorpb.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &validatorpb.SignRequest_AggregateAttestationAndProof{AggregateAttestationAndProof: agg},
	}

	byts, err := req.Marshal()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
