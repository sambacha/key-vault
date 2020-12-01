package store

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

// Paths
const (
	WalletHighestAttestationPath = "highestAttestations/"
	WalletProposalsBase          = "proposals/%s/"            // account/proposal
	WalletProposalsPath          = WalletProposalsBase + "%d" // account/proposal
)

// SaveHighestAttestation saves highest attestation
func (store *HashicorpVaultStore) SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	path := fmt.Sprintf(WalletHighestAttestationPath+"%s", store.identifierFromKey(pubKey))
	data, err := attestation.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal attestation request")
	}

	return store.storage.Put(store.ctx, &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	})
}

// RetrieveHighestAttestation retrieves highest attestation
func (store *HashicorpVaultStore) RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData {
	path := fmt.Sprintf(WalletHighestAttestationPath+"%s", store.identifierFromKey(pubKey))
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil
	}

	ret := &eth.AttestationData{}
	if err := ret.Unmarshal(entry.Value); err != nil {
		return nil
	}

	return ret
}

// SaveProposal implements Storage interface.
func (store *HashicorpVaultStore) SaveProposal(pubKey []byte, block *eth.BeaconBlock) error {
	path := fmt.Sprintf(WalletProposalsPath, store.identifierFromKey(pubKey), block.Slot)
	data, err := block.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal proposal request")
	}

	return store.storage.Put(store.ctx, &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	})
}

// RetrieveProposal implements Storage interface.
func (store *HashicorpVaultStore) RetrieveProposal(pubKey []byte, slot uint64) (*eth.BeaconBlock, error) {
	path := fmt.Sprintf(WalletProposalsPath, store.identifierFromKey(pubKey), slot)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	ret := &eth.BeaconBlock{}
	if err = ret.Unmarshal(entry.Value); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon block header object")
	}

	return ret, nil
}

// ListAllProposals returns all proposal data from the DB
func (store *HashicorpVaultStore) ListAllProposals(pubKey []byte) ([]*eth.BeaconBlock, error) {
	path := fmt.Sprintf(WalletProposalsBase, store.identifierFromKey(pubKey))
	entries, err := store.storage.List(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list records from storage with path '%s'", path)
	}

	// Return nothing if there is no record
	if len(entries) == 0 {
		return nil, nil
	}

	proposals := make([]*eth.BeaconBlock, len(entries))
	errs := make([]error, len(entries))

	var wg sync.WaitGroup
	for i, entry := range entries {
		wg.Add(1)
		go func(i int, entry string) {
			defer wg.Done()

			if entry == "latest" {
				return
			}

			epoch, err := strconv.Atoi(entry)
			if err != nil {
				errs[i] = errors.Wrapf(err, "invalid epoch number %s", entry)
				return
			}

			ret, err := store.RetrieveProposal(pubKey, uint64(epoch))
			if err != nil {
				errs[i] = errors.Wrapf(err, "failed to retrieve beacon proposal for epoch %d", epoch)
				return
			}

			proposals[i] = ret
		}(i, entry)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	clearedProposals := make([]*eth.BeaconBlock, 0)
	for _, proposal := range proposals {
		if proposal != nil {
			clearedProposals = append(clearedProposals, proposal)
		}
	}

	return clearedProposals, nil
}

func (store *HashicorpVaultStore) identifierFromKey(key []byte) string {
	return hex.EncodeToString(key)
}
