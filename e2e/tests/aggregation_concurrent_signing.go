package tests

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AggregationConcurrentSigning tests aggregation signing method concurrently.
type AggregationConcurrentSigning struct {
}

// Name returns the name of the test.
func (test *AggregationConcurrentSigning) Name() string {
	return "Test aggregation concurrent signing"
}

// Run runs the test.
func (test *AggregationConcurrentSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	store := setup.UpdateStorage(t, core.PyrmontNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, store)
	pubKey := account.ValidatorPublicKey()

	// sign and save the valid aggregation
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
	req, err := test.serializedReq(pubKey, nil, domain, agg)
	_, err = setup.SignAggregation(req, core.PyrmontNetwork)
	require.NoError(t, err)

	// Send requests in parallel
	t.Run("concurrent signing", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 5; i++ {
			i := i
			t.Run("concurrent signing "+strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()

				_, err := setup.SignAggregation(req, core.PyrmontNetwork)
				if err == nil {
					return
				}

				require.Error(t, err)
				require.IsType(t, &e2e.ServiceError{}, err)

				errValue := err.(*e2e.ServiceError).ErrorValue()
				protected := errValue == fmt.Sprintf("1 error occurred:\n\t* locked\n\n")
				require.True(t, protected, err.Error())
			})
		}
	})
}

func (test *AggregationConcurrentSigning) serializedReq(pk, root, domain []byte, agg *ethpb.AggregateAttestationAndProof) (map[string]interface{}, error) {
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
