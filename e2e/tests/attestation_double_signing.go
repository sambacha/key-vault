package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AttestationDoubleSigning tests double signing case
type AttestationDoubleSigning struct {
}

// Name returns the name of the test.
func (test *AttestationDoubleSigning) Name() string {
	return "Test double attestation signing, different block root"
}

// Run runs the test.
func (test *AttestationDoubleSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	pubKey := account.ValidatorPublicKey()

	// first sig
	att := &eth.AttestationData{
		Slot:            types.Slot(284115),
		CommitteeIndex:  types.CommitteeIndex(2),
		BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		Source: &eth.Checkpoint{
			Epoch: types.Epoch(77),
			Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		},
		Target: &eth.Checkpoint{
			Epoch: types.Epoch(78),
			Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
		},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")

	// sign and save the valid attestation
	req, err := test.serializedReq(pubKey, nil, domain, att)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)

	// second sig, different block root
	att.BeaconBlockRoot = _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0f")
	req, err = test.serializedReq(pubKey, nil, domain, att)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PraterNetwork)
	expectedErr := fmt.Sprintf("1 error occurred:\n\t* failed to sign: slashable attestation (HighestAttestationVote), not signing\n\n")
	require.Error(t, err)
	require.IsType(t, &e2e.ServiceError{}, err)
	require.EqualValues(t, expectedErr, err.(*e2e.ServiceError).ErrorValue())
}

func (test *AttestationDoubleSigning) serializedReq(pk, root, domain []byte, attestation *eth.AttestationData) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequestAttestationData{AttestationData: attestation},
	}

	byts, err := encoderv2.New().Encode(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
