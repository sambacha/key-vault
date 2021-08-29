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

	attestation := &eth.AttestationData{
		Slot:            123123,
		BeaconBlockRoot: []byte{1, 2, 3},
		CommitteeIndex:  1,
		Source: &eth.Checkpoint{
			Root:  []byte{1, 2, 3},
			Epoch: 123,
		},
		Target: &eth.Checkpoint{
			Root:  []byte{1, 2, 3},
			Epoch: 123,
		},
	}
	proposal := &eth.BeaconBlock{
		Slot:          123123,
		ProposerIndex: 1,
		ParentRoot:    []byte{1, 2, 3},
		StateRoot:     []byte{1, 2, 3},
		Body:          &eth.BeaconBlockBody{},
	}

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
		require.EqualValues(t, attestation, slashingHistory.HighestAttestation)
		require.EqualValues(t, proposal, slashingHistory.HighestProposal)
	})
}
