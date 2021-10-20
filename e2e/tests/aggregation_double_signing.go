package tests

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/prysmaticlabs/go-bitfield"

	types "github.com/prysmaticlabs/eth2-types"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AggregationDoubleSigning tests aggregation signing method concurrently.
type AggregationDoubleSigning struct {
}

// Name returns the name of the test.
func (test *AggregationDoubleSigning) Name() string {
	return "Test aggregation double signing"
}

// Run runs the test.
func (test *AggregationDoubleSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	store := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, store)
	pubKey := account.ValidatorPublicKey()

	agg := &ethpb.AggregateAttestationAndProof{
		AggregatorIndex: types.ValidatorIndex(1),
		SelectionProof:  make([]byte, 96),
		Aggregate: &ethpb.Attestation{
			AggregationBits: bitfield.NewBitlist(12),
			Signature:       make([]byte, 96),
			Data: &ethpb.AttestationData{
				Slot:            types.Slot(1),
				CommitteeIndex:  types.CommitteeIndex(12),
				BeaconBlockRoot: make([]byte, 32),
				Source: &ethpb.Checkpoint{
					Epoch: types.Epoch(1),
					Root:  make([]byte, 32),
				},
				Target: &ethpb.Checkpoint{
					Epoch: types.Epoch(1),
					Root:  make([]byte, 32),
				},
			},
		},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKey, nil, domain, agg)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)

	_, err = setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)
}

func (test *AggregationDoubleSigning) serializedReq(pk, root, domain []byte, agg *ethpb.AggregateAttestationAndProof) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequestAggregateAttestationAndProof{AggregateAttestationAndProof: agg},
	}

	byts, err := encoderv2.New().Encode(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
