package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/backend/store"
)

func TestSlashingStorage_Update(t *testing.T) {
	b, _ := getBackend(t)
	inMemStore, accountID, err := baseInmemStorage()
	require.NoError(t, err)

	account, err := inMemStore.OpenAccount(accountID)
	require.NoError(t, err)
	publicKey := hex.EncodeToString(account.ValidatorPublicKey())

	t.Run("successfully setup slashing history", func(t *testing.T) {
		slashingHistory, err := json.Marshal(struct {
			HighestAttestation *eth.AttestationData
			Proposals          []*eth.BeaconBlock `json:"proposals"`
		}{
			HighestAttestation: &eth.AttestationData{
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
			},
			Proposals: []*eth.BeaconBlock{
				{
					Slot:          123123,
					ProposerIndex: 1,
					ParentRoot:    []byte{1, 2, 3},
					StateRoot:     []byte{1, 2, 3},
					Body:          &eth.BeaconBlockBody{},
				},
			},
		})
		require.NoError(t, err)

		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		setupBaseStorage(t, req)
		req.Data = map[string]interface{}{
			publicKey: hex.EncodeToString(slashingHistory),
		}
		_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.True(t, res.Data["status"].(bool))
	})

	t.Run("rejects setup slashing history with wrong JSON format", func(t *testing.T) {
		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		setupBaseStorage(t, req)
		req.Data = map[string]interface{}{
			publicKey: hex.EncodeToString([]byte("slashinghistory")),
		}
		_, err := store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 400, res.Data["http_status_code"])
	})

	t.Run("rejects setup slashing history with wrong HEX format", func(t *testing.T) {
		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		setupBaseStorage(t, req)
		req.Data = map[string]interface{}{
			publicKey: "slashinghistory",
		}
		_, err := store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 400, res.Data["http_status_code"])
	})

	t.Run("rejects setup slashing history for unknown public key", func(t *testing.T) {
		fakePublicKey := "ab0cb36c4ce5ddabdc38a1d6868c871328539ebde5fea89686b2cd6332bf4cc5f9c48a501d1d6d87bf916d0e1b01ead963e1b6ce52075e26dc65bad535ecfad0"
		slashingHistory, err := json.Marshal(struct {
			Attestations []*eth.AttestationData `json:"attestations"`
			Proposals    []*eth.BeaconBlock     `json:"proposals"`
		}{
			Attestations: []*eth.AttestationData{},
			Proposals:    []*eth.BeaconBlock{},
		})
		require.NoError(t, err)

		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		setupBaseStorage(t, req)
		req.Data = map[string]interface{}{
			fakePublicKey: hex.EncodeToString(slashingHistory),
		}
		_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 404, res.Data["http_status_code"])
	})
}

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
		err = newStore.SaveProposal(account.ValidatorPublicKey(), proposal)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, res.Data[publicKey].(string))

		data, err := hex.DecodeString(res.Data[publicKey].(string))
		require.NoError(t, err)
		var slashingHistory SlashingHistory
		err = json.Unmarshal(data, &slashingHistory)
		require.NoError(t, err)
		//require.Len(t, slashingHistory.Attestations, 1)
		require.Len(t, slashingHistory.Proposals, 1)
		require.EqualValues(t, attestation, slashingHistory.HighestAttestation)
		require.EqualValues(t, proposal, slashingHistory.Proposals[0])
	})
}
