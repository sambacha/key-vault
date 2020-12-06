package store

import (
	"encoding/hex"
	"fmt"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

// Paths
const (
	WalletHighestAttestationPath = "highestAttestations/"
	WalletHighestProposalsBase   = "proposals/%s" // account/proposal
)

// SaveHighestAttestation saves highest attestation
func (store *HashicorpVaultStore) SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	if attestation == nil || pubKey == nil {
		return errors.Errorf("pubKey and attestation must not be nil")
	}

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
	if pubKey == nil {
		return nil
	}

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

// SaveHighestProposal implements Storage interface.
func (store *HashicorpVaultStore) SaveHighestProposal(pubKey []byte, block *eth.BeaconBlock) error {
	if block == nil || pubKey == nil {
		return errors.Errorf("pubKey and block must not be nil")
	}

	path := fmt.Sprintf(WalletHighestProposalsBase, store.identifierFromKey(pubKey))
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

// RetrieveHighestProposal implements Storage interface.
func (store *HashicorpVaultStore) RetrieveHighestProposal(pubKey []byte) *eth.BeaconBlock {
	if pubKey == nil {
		return nil
	}

	path := fmt.Sprintf(WalletHighestProposalsBase, store.identifierFromKey(pubKey))
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil
	}

	ret := &eth.BeaconBlock{}
	if err = ret.Unmarshal(entry.Value); err != nil {
		return nil
	}

	return ret
}

func (store *HashicorpVaultStore) identifierFromKey(key []byte) string {
	return hex.EncodeToString(key)
}
