package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	types "github.com/prysmaticlabs/eth2-types"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AttestationFarFutureSigning tests sign attestation endpoint with future signing.
type AttestationFarFutureSigning struct {
}

// Name returns the name of the test.
func (test *AttestationFarFutureSigning) Name() string {
	return "Test far future attestation (source and target) signing"
}

// Run run the test.
func (test *AttestationFarFutureSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	expectedSourceErr := fmt.Sprintf("map[string]interface {}{\"errors\":[]interface {}{\"1 error occurred:\\n\\t* failed to sign: source epoch too far into the future\\n\\n\"}}")
	expectedTargetErr := fmt.Sprintf("map[string]interface {}{\"errors\":[]interface {}{\"1 error occurred:\\n\\t* failed to sign: target epoch too far into the future\\n\\n\"}}")

	currentEpoch := core.PraterNetwork.EstimatedCurrentEpoch()
	futureEpoch := core.PraterNetwork.EstimatedCurrentEpoch() + 1000
	test.testFarFuture(t, setup, pubKeyBytes, futureEpoch, currentEpoch, expectedSourceErr) // far future source
	test.testFarFuture(t, setup, pubKeyBytes, currentEpoch, futureEpoch, expectedTargetErr) // far future target
	test.testFarFuture(t, setup, pubKeyBytes, futureEpoch, futureEpoch, expectedTargetErr)  // far future both
}

func (test *AttestationFarFutureSigning) testFarFuture(
	t *testing.T,
	setup *e2e.BaseSetup,
	pubKeyBytes []byte,
	source types.Epoch,
	target types.Epoch,
	expectedErr string,
) {
	att := &eth.AttestationData{
		Slot:            core.PraterNetwork.EstimatedCurrentSlot() + 1000,
		CommitteeIndex:  types.CommitteeIndex(2),
		BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		Source: &eth.Checkpoint{
			Epoch: source,
			Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		},
		Target: &eth.Checkpoint{
			Epoch: target,
			Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
		},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")

	// Send sign attestation request
	req, err := test.serializedReq(pubKeyBytes, nil, domain, att)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PraterNetwork)
	require.NotNil(t, err)
	require.EqualError(t, err, expectedErr, fmt.Sprintf("actual: %s\n", err.Error()))
}

func (test *AttestationFarFutureSigning) serializedReq(pk, root, domain []byte, attestation *eth.AttestationData) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequest_AttestationData{AttestationData: attestation},
	}

	byts, err := encoderv2.New().Encode(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
