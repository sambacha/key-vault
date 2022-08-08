package tests

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	"github.com/bloxapp/eth2-key-manager/signer"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
	slashingprotection "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
	consensusblocks "github.com/prysmaticlabs/prysm/consensus-types/blocks"
)

// ProposalSigning tests sign proposal endpoint.
type ProposalSigning struct {
}

// Name returns the name of the test.
func (test *ProposalSigning) Name() string {
	return "Test proposal signing"
}

// Run run the test.
func (test *ProposalSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.HDWallet, nil)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	// Get wallet
	wallet, err := storage.OpenWallet()
	require.NoError(t, err)

	blk := referenceBlock(t)
	domain := _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
	req, err := test.serializedReq(pubKeyBytes, nil, domain, blk)
	require.NoError(t, err)

	// Sign data
	protector := slashingprotection.NewNormalProtection(storage)
	var signer signer.ValidatorSigner = signer.NewSimpleSigner(wallet, protector, storage.Network())

	wrappedBlk, err := consensusblocks.NewBeaconBlock(blk)
	require.NoError(t, err)
	res, err := signer.SignBeaconBlock(wrappedBlk, domain, pubKeyBytes)
	require.NoError(t, err)

	// Send sign attestation request
	sig, err := setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)

	require.Equal(t, res, sig)
}

func (test *ProposalSigning) serializedReq(pk, root, domain []byte, blk *eth.BeaconBlock) (map[string]interface{}, error) {
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
