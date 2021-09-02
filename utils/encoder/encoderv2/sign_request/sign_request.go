package sign_request

import (
	types "github.com/prysmaticlabs/eth2-types"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

type ISignObject interface {
	isSignRequest_Object()
}

type SignRequest struct {
	PublicKey       []byte      `json:"public_key,omitempty"`
	SigningRoot     []byte      `json:"signing_root,omitempty"`
	SignatureDomain []byte      `json:"signature_domain,omitempty"`
	Object          ISignObject `json:"object,omitempty"`
}

func (x *SignRequest) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *SignRequest) GetSigningRoot() []byte {
	if x != nil {
		return x.SigningRoot
	}
	return nil
}

func (x *SignRequest) GetSignatureDomain() []byte {
	if x != nil {
		return x.SignatureDomain
	}
	return nil
}

func (m *SignRequest) GetObject() ISignObject {
	if m != nil {
		return m.Object
	}
	return nil
}

func (x *SignRequest) GetBlock() *eth.BeaconBlock {
	if x, ok := x.GetObject().(*SignRequest_Block); ok {
		return x.Block
	}
	return nil
}

func (x *SignRequest) GetAttestationData() *eth.AttestationData {
	if x, ok := x.GetObject().(*SignRequest_AttestationData); ok {
		return x.AttestationData
	}
	return nil
}

func (x *SignRequest) GetAggregateAttestationAndProof() *eth.AggregateAttestationAndProof {
	if x, ok := x.GetObject().(*SignRequest_AggregateAttestationAndProof); ok {
		return x.AggregateAttestationAndProof
	}
	return nil
}

func (x *SignRequest) GetExit() *eth.VoluntaryExit {
	if x, ok := x.GetObject().(*SignRequest_Exit); ok {
		return x.Exit
	}
	return nil
}

func (x *SignRequest) GetSlot() types.Slot {
	if x, ok := x.GetObject().(*SignRequest_Slot); ok {
		return x.Slot
	}
	return types.Slot(0)
}

func (x *SignRequest) GetEpoch() types.Epoch {
	if x, ok := x.GetObject().(*SignRequest_Epoch); ok {
		return x.Epoch
	}
	return types.Epoch(0)
}

func (x *SignRequest) GetBlockV2() *eth.BeaconBlockAltair {
	if x, ok := x.GetObject().(*SignRequest_BlockV2); ok {
		return x.BlockV2
	}
	return nil
}

func (x *SignRequest) GetSyncAggregatorSelectionData() *eth.SyncAggregatorSelectionData {
	if x, ok := x.GetObject().(*SignRequest_SyncAggregatorSelectionData); ok {
		return x.SyncAggregatorSelectionData
	}
	return nil
}

func (x *SignRequest) GetContributionAndProof() *eth.ContributionAndProof {
	if x, ok := x.GetObject().(*SignRequest_ContributionAndProof); ok {
		return x.ContributionAndProof
	}
	return nil
}

func (x *SignRequest) GetSyncCommitteeMessage() types.SSZBytes {
	if x, ok := x.GetObject().(*SignRequest_SyncCommitteeMessage); ok {
		return x.Root
	}
	return nil
}
