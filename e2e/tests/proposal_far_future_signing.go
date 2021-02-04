package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// ProposalFarFutureSigning tests sign proposal endpoint with future signing.
type ProposalFarFutureSigning struct {
}

// Name returns the name of the test.
func (test *ProposalFarFutureSigning) Name() string {
	return "Test proposal far future signing"
}

// Run run the test.
func (test *ProposalFarFutureSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PyrmontNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	blk := referenceBlock(t)
	blk.Slot = core.PyrmontNetwork.EstimatedCurrentSlot() + 100
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKeyBytes, nil, domain, blk)
	require.NoError(t, err)
	_, err = setup.Sign("sign", req, core.PyrmontNetwork)
	require.NotNil(t, err)
	expectedErr := fmt.Sprintf("map[string]interface {}{\"errors\":[]interface {}{\"1 error occurred:\\n\\t* failed to sign: proposed block slot too far into the future\\n\\n\"}}")
	require.EqualError(t, err, expectedErr, fmt.Sprintf("actual: %s\n", err.Error()))
}

func (test *ProposalFarFutureSigning) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
	req := &validatorpb.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &validatorpb.SignRequest_Block{Block: blk},
	}

	byts, err := req.Marshal()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
