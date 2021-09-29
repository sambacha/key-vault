package shared

import (
	"encoding/hex"
	"fmt"
	"testing"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
	"github.com/bloxapp/eth2-key-manager/wallets"
	"github.com/bloxapp/eth2-key-manager/wallets/hd"
	"github.com/bloxapp/eth2-key-manager/wallets/nd"
)

// AccountIndex is the test account index.
const AccountIndex = 0

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

// BaseInmemStorage creates the in-memory storage and creates the base account.
func BaseInmemStorage(t *testing.T, minimalSlashingData bool, walletType core.WalletType, privKey []byte) (*inmemory.InMemStore, error) {
	err := core.InitBLS()
	require.NoError(t, err)

	store := inmemory.NewInMemStore(core.PraterNetwork)

	// wallet
	walletCtx := &core.WalletContext{Storage: store}

	// account
	var acc core.ValidatorAccount
	if walletType == core.NDWallet {
		wallet := nd.NewWallet(walletCtx)
		if err := store.SaveWallet(wallet); err != nil {
			return nil, err
		}
		k, err := core.NewHDKeyFromPrivateKey(privKey, "")
		require.NoError(t, err)
		acc = wallets.NewValidatorAccount("", k, k.PublicKey().Serialize(), "", walletCtx)
		require.NoError(t, wallet.AddValidatorAccount(acc))
	} else {
		wallet := hd.NewWallet(walletCtx)
		if err := store.SaveWallet(wallet); err != nil {
			return nil, err
		}
		acc, err = wallet.CreateValidatorAccount(newSeed(t), nil)
		if err != nil {
			return nil, err
		}
		err = store.SaveAccount(acc)
		if err != nil {
			return nil, err
		}
	}

	// base highest att.
	if minimalSlashingData {
		err = store.SaveHighestAttestation(acc.ValidatorPublicKey(), &eth.AttestationData{
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  nil,
			},
			Target: &eth.Checkpoint{
				Epoch: 0,
				Root:  nil,
			},
		})
		if err != nil {
			return nil, err
		}

		err = store.SaveHighestProposal(acc.ValidatorPublicKey(), &eth.BeaconBlock{
			Slot: 0,
		})
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

// RetrieveAccount retrieves test account fro the storage.
func RetrieveAccount(t *testing.T, store core.Storage) core.ValidatorAccount {
	accounts, err := store.ListAccounts()
	require.NoError(t, err)

	for _, acc := range accounts {
		if acc.BasePath() == fmt.Sprintf("/%d", AccountIndex) {
			return acc
		}
	}
	return nil
}

func newSeed(t *testing.T) []byte {
	entropy, err := core.GenerateNewEntropy()
	require.NoError(t, err)

	seed, err := core.SeedFromEntropy(entropy, "")
	require.NoError(t, err)
	return seed
}
