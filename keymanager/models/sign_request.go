package models

import (
	types "github.com/prysmaticlabs/eth2-types"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// ISignObject interface
type ISignObject interface {
	isSignRequestObject()
}

// SignRequest implementing ISignObject
type SignRequest struct {
	PublicKey       []byte      `json:"public_key,omitempty"`
	SigningRoot     []byte      `json:"signing_root,omitempty"`
	SignatureDomain []byte      `json:"signature_domain,omitempty"`
	Object          ISignObject `json:"object,omitempty"`
}

// GetPublicKey return publicKey
func (x *SignRequest) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

// GetSigningRoot return root bytes
func (x *SignRequest) GetSigningRoot() []byte {
	if x != nil {
		return x.SigningRoot
	}
	return nil
}

// GetSignatureDomain return domain bytes
func (x *SignRequest) GetSignatureDomain() []byte {
	if x != nil {
		return x.SignatureDomain
	}
	return nil
}

// GetObject return ISignObject interface
func (x *SignRequest) GetObject() ISignObject {
	if x != nil {
		return x.Object
	}
	return nil
}

// GetBlock return req block
func (x *SignRequest) GetBlock() *eth.BeaconBlock {
	if x, ok := x.GetObject().(*SignRequestBlock); ok {
		return x.Block
	}
	return nil
}

// GetAttestationData return AttestationData
func (x *SignRequest) GetAttestationData() *eth.AttestationData {
	if x, ok := x.GetObject().(*SignRequestAttestationData); ok {
		return x.AttestationData
	}
	return nil
}

// GetAggregateAttestationAndProof return AggregateAttestationAndProof
func (x *SignRequest) GetAggregateAttestationAndProof() *eth.AggregateAttestationAndProof {
	if x, ok := x.GetObject().(*SignRequestAggregateAttestationAndProof); ok {
		return x.AggregateAttestationAndProof
	}
	return nil
}

// GetExit return VoluntaryExit
func (x *SignRequest) GetExit() *eth.VoluntaryExit {
	if x, ok := x.GetObject().(*SignRequestExit); ok {
		return x.Exit
	}
	return nil
}

// GetSlot return types slot
func (x *SignRequest) GetSlot() types.Slot {
	if x, ok := x.GetObject().(*SignRequestSlot); ok {
		return x.Slot
	}
	return types.Slot(0)
}

// GetEpoch return types epoch
func (x *SignRequest) GetEpoch() types.Epoch {
	if x, ok := x.GetObject().(*SignRequestEpoch); ok {
		return x.Epoch
	}
	return types.Epoch(0)
}

// GetBlockV2 return altair block
func (x *SignRequest) GetBlockV2() *eth.BeaconBlockAltair {
	if x, ok := x.GetObject().(*SignRequestBlockV2); ok {
		return x.BlockV2
	}
	return nil
}

// GetSyncAggregatorSelectionData return SyncAggregatorSelectionData
func (x *SignRequest) GetSyncAggregatorSelectionData() *eth.SyncAggregatorSelectionData {
	if x, ok := x.GetObject().(*SignRequestSyncAggregatorSelectionData); ok {
		return x.SyncAggregatorSelectionData
	}
	return nil
}

// GetContributionAndProof return ContributionAndProof
func (x *SignRequest) GetContributionAndProof() *eth.ContributionAndProof {
	if x, ok := x.GetObject().(*SignRequestContributionAndProof); ok {
		return x.ContributionAndProof
	}
	return nil
}

// GetSyncCommitteeMessage return types SSZBytes
func (x *SignRequest) GetSyncCommitteeMessage() types.SSZBytes {
	if x, ok := x.GetObject().(*SignRequestSyncCommitteeMessage); ok {
		return x.Root
	}
	return nil
}
