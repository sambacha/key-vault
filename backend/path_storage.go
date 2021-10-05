package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/backend/store"
)

// Endpoints patterns
const (
	// StoragePattern is the path pattern for storage endpoint
	StoragePattern = "storage"
)

func storagePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         StoragePattern,
			HelpSynopsis:    "Update storage using one or more accounts",
			HelpDescription: `Manage KeyVault storage`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "storage to update",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathStorageUpdateV2,
			},
		},
	}
}

// buildInMemoryStorage returns unmarshalled in-memory storage
func buildInMemStore(data *framework.FieldData) (*inmemory.InMemStore, error) {
	storage := data.Get("data").(string)
	storageBytes, err := hex.DecodeString(storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode storage")
	}

	var inMemStore *inmemory.InMemStore
	err = json.Unmarshal(storageBytes, &inMemStore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON un-marshal storage")
	}
	return inMemStore, nil
}

// pathStorageUpdate updates all accounts from new uploaded storage
func (b *backend) pathStorageUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	inMemStore, err := buildInMemStore(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build in memory store")
	}

	_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}

// pathStorageUpdateV2 updates storage accounts from new requested storage
func (b *backend) pathStorageUpdateV2(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	inMemStore, err := buildInMemStore(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build in memory store")
	}

	// Update hashicorp store with new account(s)
	_, err = store.FromInMemoryStoreV2(ctx, inMemStore, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage from in memory")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}
