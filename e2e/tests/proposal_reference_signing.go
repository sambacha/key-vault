package tests

import (
	"encoding/hex"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/utils/encoder"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/keymanager/models"
)

func referenceBlock(t *testing.T) *spec.VersionedBeaconBlock {
	blockByts := _byteArray("020000000000000055000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776badd5cb7e6a4bffd8ce7fe9697aed511661861e312ad546dcf5480159698f47a554000000a2c156a4bc9439f1d85f922f2abaa96e830f1c526101211bdb7d16f4ad9490a0302fc5adb089c05b5f16fd465962f47c04fc2b81a94d135a07c1613db61511c17284b51fafab984e56d3411e16e45f5068f146d9412f91d31ab0f237eac3d745a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a8000000000000000a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a0000000000000000000000000000000000000000000000000000000000000000dc000000dc000000dc000000c5010000c501000004000000e4000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b97b6f271ac364b041cd465f32fa7ffa19f5a811f1e6e14713f93e06537ef827d382bac72f0990b84f83cd9bbe0062815020086bf27b9ced172cc6add8ba5197991cf634d18666f5d43df6f09180ce20a357e4d05b2784409e32147f1042986e31f")
	blk := &phase0.BeaconBlock{}
	require.NoError(t, blk.UnmarshalSSZ(blockByts))
	return &spec.VersionedBeaconBlock{
		Phase0: blk,
	}
}

// ProposalReferenceSigning tests sign proposal endpoint.
type ProposalReferenceSigning struct {
}

// Name returns the name of the test.
// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/master/shared/testutil/block.go#L86
func (test *ProposalReferenceSigning) Name() string {
	return "Test proposal reference signing"
}

// Run run the test.
func (test *ProposalReferenceSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PraterNetwork, true, core.NDWallet, _byteArray("5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc"))
	// Get wallet
	wallet, err := storage.OpenWallet()
	require.NoError(t, err)
	account, err := wallet.AccountByPublicKey("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972")
	require.NoError(t, err)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	blk := referenceBlock(t)
	domain := _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
	req, err := test.serializedReq(pubKeyBytes, nil, domain, blk)
	require.NoError(t, err)
	sig, err := setup.Sign("sign", req, core.PraterNetwork)
	require.NoError(t, err)

	expectedSig := _byteArray("83ed05e039d39971dab6bd224809df2f646425526f6bdbc0b0ac3f76ab73585397f2bd85b77f34503a6bc404f9a651cc0b4711e4d1a8596e3195e8c64df1f30caf3be1067251333221888dd1569147ce3ddd77787feb99c43d9e25cd64bd755b")
	require.Equal(t, expectedSig, sig)
}

func (test *ProposalReferenceSigning) serializedReq(pk, root []byte, domain [32]byte, blk *spec.VersionedBeaconBlock) (map[string]interface{}, error) {
	req := &models.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &models.SignRequestBlock{VersionedBeaconBlock: blk},
	}

	byts, err := encoder.New().Encode(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
