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
	StoragePattern         = "storage"
	SeedlessStoragePattern = "storage/seedless"
)

func storagePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         StoragePattern,
			HelpSynopsis:    "Update storage",
			HelpDescription: `Manage KeyVault storage`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "storage to update",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathStorageUpdate,
			},
		},
		{
			Pattern:         SeedlessStoragePattern,
			HelpSynopsis:    "Update storage using one account",
			HelpDescription: `Manage KeyVault storage in a seedless mode`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "storage to update",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathStorageUpdateSeedless,
			},
		},
	}
}

// getVaultMemoryStorage returns unmarshalled in-memory storage
func getVaultMemoryStorage(data *framework.FieldData) (*inmemory.InMemStore, error) {
	storage := data.Get("data").(string)
	storageBytes, err := hex.DecodeString(storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode storage")
	}

	var vaultMemoryStorage *inmemory.InMemStore
	err = json.Unmarshal(storageBytes, &vaultMemoryStorage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON un-marshal storage")
	}
	return vaultMemoryStorage, nil
}

// pathStorageUpdate updates all accounts from new uploaded storage
func (b *backend) pathStorageUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vaultMemoryStorage, err := getVaultMemoryStorage(data)
	if err != nil {
		return nil, err
	}

	_, err = store.FromInMemoryStore(ctx, vaultMemoryStorage, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}

// pathStorageUpdateSeedless updates accounts from one account storage
func (b *backend) pathStorageUpdateSeedless(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vaultMemoryStorage, err := getVaultMemoryStorage(data)
	if err != nil {
		return nil, err
	}

	_, err = store.FromInMemoryStoreSeedless(ctx, vaultMemoryStorage, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}
