package tests

import (
	"encoding/hex"
	"net/http"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

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
	setup.UpdateStorage(t, core.PyrmontNetwork, true, core.HDWallet, nil)

	// sign
	blk := &eth.BeaconBlock{
		Slot:          78,
		ProposerIndex: 1010,
		ParentRoot:    _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		StateRoot:     _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		Body:          &eth.BeaconBlockBody{},
	}
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(make([]byte, 48), nil, domain, blk)
	require.NoError(t, err)

	_, err = setup.SignProposal(req, core.PyrmontNetwork)
	require.Error(t, err)
	require.IsType(t, &e2e.ServiceError{}, err)
	require.EqualValues(t, "account not found", err.(*e2e.ServiceError).DataValue("message"))
	require.EqualValues(t, http.StatusNotFound, err.(*e2e.ServiceError).DataValue("status_code"))
}

func (test *ProposalSigningAccountNotFound) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
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
