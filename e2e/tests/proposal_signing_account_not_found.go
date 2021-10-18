package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
)

// ProposalSigningAccountNotFound tests sign attestation when account not found
type ProposalSigningAccountNotFound struct {
}

// Name returns the name of the test
func (test *ProposalSigningAccountNotFound) Name() string {
	return "Test proposal signing account not found"
}

// Run runs the test.
func (test *ProposalSigningAccountNotFound) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)

	// sign
	blk := referenceBlock(t)
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(make([]byte, 48), nil, domain, blk)
	require.NoError(t, err)

	_, err = setup.Sign("sign", req, core.PraterNetwork)
	require.Error(t, err)
	require.IsType(t, &e2e.ServiceError{}, err)
	require.EqualValues(t, fmt.Sprintf("1 error occurred:\n\t* failed to sign: account not found\n\n"), err.(*e2e.ServiceError).ErrorValue())
}

func (test *ProposalSigningAccountNotFound) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequestBlock{Block: blk},
	}

	byts, err := encoderv2.New().Encode(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
