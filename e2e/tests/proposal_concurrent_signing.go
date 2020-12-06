package tests

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// ProposalConcurrentSigning tests signing method concurrently.
type ProposalConcurrentSigning struct {
}

// Name returns the name of the test.
func (test *ProposalConcurrentSigning) Name() string {
	return "Test proposal concurrent signing"
}

// Run runs the test.
func (test *ProposalConcurrentSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	store := setup.UpdateStorage(t, core.PyrmontNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, store)
	pubKey := account.ValidatorPublicKey()

	// sign and save the valid proposal
	blk := referenceBlock(t)
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKey, nil, domain, blk)
	require.NoError(t, err)

	_, err = setup.Sign("sign", req, core.PyrmontNetwork)
	require.NoError(t, err)

	// Send requests in parallel
	t.Run("concurrent signing", func(t *testing.T) {
		t.Parallel()
		for i := uint64(2); i < 7; i++ {
			i := i
			t.Run("concurrent signing "+strconv.Itoa(int(i)), func(t *testing.T) {
				t.Parallel()
				test.runSlashableProposal(t, setup, pubKey, i)
			})
		}
	})
}

// will return no error if trying to sign a slashable attestation will not work
func (test *ProposalConcurrentSigning) runSlashableProposal(t *testing.T, setup *e2e.BaseSetup, pubKey []byte, index uint64) {
	blk := referenceBlock(t)
	blk.ParentRoot = _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0d")
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKey, nil, domain, blk)
	require.NoError(t, err)

	_, err = setup.Sign("sign", req, core.PyrmontNetwork)
	require.Error(t, err, "did not slash", index)
	require.IsType(t, &e2e.ServiceError{}, err)

	errValue := err.(*e2e.ServiceError).ErrorValue()
	protected := errValue == fmt.Sprintf("1 error occurred:\n\t* failed to sign: slashable proposal (HighestProposalVote), not signing\n\n") ||
		errValue == fmt.Sprintf("1 error occurred:\n\t* locked\n\n")
	require.True(t, protected, err.Error())
}

func (test *ProposalConcurrentSigning) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
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
