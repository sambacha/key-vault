package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/backend/store"
)

func TestSlashingStorage_Read(t *testing.T) {
	b, _ := getBackend(t)
	inMemStore, accountID, err := baseInmemStorage()
	require.NoError(t, err)

	account, err := inMemStore.OpenAccount(accountID)
	require.NoError(t, err)
	publicKey := hex.EncodeToString(account.ValidatorPublicKey())

	attestationDataSSZ := _byteArray("000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b")
	blockSSZ := _byteArray("010000000000000055000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776badd5cb7e6a4bffd8ce7fe9697aed511661861e312ad546dcf5480159698f47a554000000a2c156a4bc9439f1d85f922f2abaa96e830f1c526101211bdb7d16f4ad9490a0302fc5adb089c05b5f16fd465962f47c04fc2b81a94d135a07c1613db61511c17284b51fafab984e56d3411e16e45f5068f146d9412f91d31ab0f237eac3d745a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a8000000000000000a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a0000000000000000000000000000000000000000000000000000000000000000dc000000dc000000dc000000c5010000c501000004000000e4000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b97b6f271ac364b041cd465f32fa7ffa19f5a811f1e6e14713f93e06537ef827d382bac72f0990b84f83cd9bbe0062815020086bf27b9ced172cc6add8ba5197991cf634d18666f5d43df6f09180ce20a357e4d05b2784409e32147f1042986e31f")

	attestation := &eth.AttestationData{}
	require.NoError(t, attestation.UnmarshalSSZ(attestationDataSSZ))

	proposal := &eth.BeaconBlock{}
	require.NoError(t, proposal.UnmarshalSSZ(blockSSZ))

	t.Run("successfully read slashing history", func(t *testing.T) {
		ctx := context.Background()
		req := logical.TestRequest(t, logical.ReadOperation, "storage/slashing")
		setupBaseStorage(t, req)
		newStore, err := store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)
		err = newStore.SaveHighestAttestation(account.ValidatorPublicKey(), attestation)
		require.NoError(t, err)
		err = newStore.SaveHighestProposal(account.ValidatorPublicKey(), proposal)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, res.Data[publicKey].(string))

		data, err := hex.DecodeString(res.Data[publicKey].(string))
		require.NoError(t, err)
		var slashingHistory SlashingHistory
		err = json.Unmarshal(data, &slashingHistory)
		require.NoError(t, err)
		decodedAttSSZ, err := slashingHistory.HighestAttestation.MarshalSSZ()
		require.NoError(t, err)
		decodedProposalSSZ, err := slashingHistory.HighestProposal.MarshalSSZ()
		require.NoError(t, err)
		require.EqualValues(t, attestationDataSSZ, decodedAttSSZ)
		require.EqualValues(t, blockSSZ, decodedProposalSSZ)
	})
}
