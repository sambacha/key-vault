package store

import (
	"context"
	"encoding/json"
	"fmt"

	vault "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
	"github.com/bloxapp/eth2-key-manager/wallets"
	"github.com/bloxapp/eth2-key-manager/wallets/hd"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/utils/encoder"
)

// Paths
const (
	WalletDataPath = "wallet/data"
	AccountBase    = "wallet/accounts/"
	AccountPath    = AccountBase + "%s"
)

// HashicorpVaultStore implements store.Store interface using Vault.
type HashicorpVaultStore struct {
	storage logical.Storage
	ctx     context.Context
	network core.Network
	encoder encoder.IEncoder

	encryptor          encryptor.Encryptor
	encryptionPassword []byte
}

// NewHashicorpVaultStore is the constructor of HashicorpVaultStore.
func NewHashicorpVaultStore(ctx context.Context, storage logical.Storage, network core.Network) *HashicorpVaultStore {
	return &HashicorpVaultStore{
		storage: storage,
		network: network,
		ctx:     ctx,
		encoder: encoder.New(),
	}
}

// FromInMemoryStoreV2 updates HashicorpVaultStore with new accounts.
func FromInMemoryStoreV2(ctx context.Context, newStorage *inmemory.InMemStore, existingStorage logical.Storage) (*HashicorpVaultStore, error) {

	// Open newStorage wallet
	newStorageWallet, err := newStorage.OpenWallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open newStorage wallet")
	}

	// Get existing hashicorp storage
	hashicorpStore := NewHashicorpVaultStore(ctx, existingStorage, newStorage.Network())
	options := vault.KeyVaultOptions{}
	options.SetStorage(hashicorpStore)

	_, err = vault.OpenKeyVault(&options)
	// If no existing hashicorp store - use in memory store
	if err != nil {
		// Save wallet in hashicorp store
		err = hashicorpStore.SaveWallet(newStorageWallet)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save wallet to hashicorp store")
		}
	}

	// Open existing wallet
	existingWallet, err := hashicorpStore.OpenWallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open existing wallet")
	}

	// Save new accounts
	for _, newAccount := range newStorageWallet.Accounts() {
		// Check if account already exists and don't change it
		// TODO: how to handle the same account name with index but different public keys?
		existingAccount, _ := existingWallet.AccountByPublicKey(string(newAccount.ValidatorPublicKey()))
		if existingAccount != nil {
			continue
		}

		// Add validator account in wallet
		err := existingWallet.AddValidatorAccount(newAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save account")
		}

		// Save account in vault
		if err := hashicorpStore.SaveAccount(newAccount); err != nil {
			return nil, errors.Wrap(err, "failed to save account")
		}

		// Save highest attestation
		highestAtt, found, err := newStorage.RetrieveHighestAttestation(newAccount.ValidatorPublicKey())
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve highest attestation")
		}
		if found && highestAtt != nil {
			if err := hashicorpStore.SaveHighestAttestation(newAccount.ValidatorPublicKey(), highestAtt); err != nil {
				return nil, errors.Wrap(err, "failed to save highest attestation")
			}
		}

		// Save highest proposal
		highestProposal, found, err := newStorage.RetrieveHighestProposal(newAccount.ValidatorPublicKey())
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve highest attestation")
		}
		if found && highestProposal != 0 {
			if err := hashicorpStore.SaveHighestProposal(newAccount.ValidatorPublicKey(), highestProposal); err != nil {
				return nil, errors.Wrap(err, "failed to save highest proposal")
			}
		}
	}

	return hashicorpStore, nil
}

// FromInMemoryStore creates the HashicorpVaultStore based on the given in-memory store.
func FromInMemoryStore(ctx context.Context, newStorage *inmemory.InMemStore, existingStorage logical.Storage) (*HashicorpVaultStore, error) {
	// first delete old data
	// delete all accounts
	res, err := existingStorage.List(ctx, AccountBase)
	if err != nil {
		return nil, err
	}

	for _, accountID := range res {
		path := fmt.Sprintf(AccountPath, accountID)
		if err := existingStorage.Delete(ctx, path); err != nil {
			return nil, err
		}
	}

	if err := existingStorage.Delete(ctx, WalletDataPath); err != nil {
		return nil, err
	}

	if err := existingStorage.Delete(ctx, AccountBase); err != nil {
		return nil, err
	}

	if err := existingStorage.Delete(ctx, WalletHighestAttestationPath); err != nil {
		return nil, err
	}

	if err := existingStorage.Delete(ctx, WalletHighestProposalsBase); err != nil {
		return nil, err
	}

	// Create new store
	newHashicorpVaultStore := NewHashicorpVaultStore(ctx, existingStorage, newStorage.Network())

	// Save wallet
	wallet, err := newStorage.OpenWallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open wallet")
	}

	if err := newHashicorpVaultStore.SaveWallet(wallet); err != nil {
		return nil, errors.Wrap(err, "failed to save wallet")
	}

	// Save accounts
	for _, acc := range wallet.Accounts() {
		if err := newHashicorpVaultStore.SaveAccount(acc); err != nil {
			return nil, errors.Wrap(err, "failed to save account")
		}
	}

	// save highest att.
	for _, acc := range wallet.Accounts() {
		highestAtt, found, err := newStorage.RetrieveHighestAttestation(acc.ValidatorPublicKey())
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve highest attestation")
		}
		if found && highestAtt != nil {
			if err := newHashicorpVaultStore.SaveHighestAttestation(acc.ValidatorPublicKey(), highestAtt); err != nil {
				return nil, errors.Wrap(err, "failed to save highest attestation")
			}
		}
	}

	// save highest proposal.
	for _, acc := range wallet.Accounts() {
		// Save highest proposal
		highestProposal, found, err := newStorage.RetrieveHighestProposal(acc.ValidatorPublicKey())
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve highest attestation")
		}
		if found && highestProposal != 0 {
			if err := newHashicorpVaultStore.SaveHighestProposal(acc.ValidatorPublicKey(), highestProposal); err != nil {
				return nil, errors.Wrap(err, "failed to save highest proposal")
			}
		}
	}

	return newHashicorpVaultStore, nil
}

// Name returns the name of the store.
func (store *HashicorpVaultStore) Name() string {
	return "Hashicorp Vault"
}

// Network returns the network the storage is related to.
func (store *HashicorpVaultStore) Network() core.Network {
	return store.network
}

// SaveWallet implements Storage interface.
func (store *HashicorpVaultStore) SaveWallet(wallet core.Wallet) error {
	data, err := json.Marshal(wallet)
	if err != nil {
		return errors.Wrap(err, "failed to marshal wallet")
	}

	return store.storage.Put(store.ctx, &logical.StorageEntry{
		Key:      WalletDataPath,
		Value:    data,
		SealWrap: false,
	})
}

// OpenWallet returns nil,nil if no wallet was found
func (store *HashicorpVaultStore) OpenWallet() (core.Wallet, error) {
	entry, err := store.storage.Get(store.ctx, WalletDataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet data")
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, errors.New("wallet not found")
	}

	var ret hd.Wallet
	ret.SetContext(store.freshContext())
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal HD Wallet object")
	}

	return &ret, nil
}

// ListAccounts returns an empty array for no accounts
func (store *HashicorpVaultStore) ListAccounts() ([]core.ValidatorAccount, error) {
	w, err := store.OpenWallet()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get wallet")
	}

	return w.Accounts(), nil
}

// SaveAccount stores the given account in DB.
func (store *HashicorpVaultStore) SaveAccount(account core.ValidatorAccount) error {
	data, err := json.Marshal(account)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account object")
	}

	return store.storage.Put(store.ctx, &logical.StorageEntry{
		Key:      fmt.Sprintf(AccountPath, account.ID().String()),
		Value:    data,
		SealWrap: false,
	})
}

// OpenAccount opens an account by the given ID. Returns nil,nil if no account was found.
func (store *HashicorpVaultStore) OpenAccount(accountID uuid.UUID) (core.ValidatorAccount, error) {
	path := fmt.Sprintf(AccountPath, accountID)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	var ret wallets.HDAccount
	ret.SetContext(store.freshContext())
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal HD account object")
	}
	return &ret, nil
}

// DeleteAccount deletes the given account
func (store *HashicorpVaultStore) DeleteAccount(accountID uuid.UUID) error {
	path := fmt.Sprintf(AccountPath, accountID)
	if err := store.storage.Delete(store.ctx, path); err != nil {
		return errors.Wrapf(err, "failed to delete record with path '%s'", path)
	}
	return nil
}

// SetEncryptor sets the given encryptor. Could be nil value.
func (store *HashicorpVaultStore) SetEncryptor(encryptor encryptor.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *HashicorpVaultStore) freshContext() *core.WalletContext {
	return &core.WalletContext{
		Storage: store,
	}
}
