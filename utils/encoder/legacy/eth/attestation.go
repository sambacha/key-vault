package eth

import (
	types "github.com/prysmaticlabs/eth2-types"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func NewAttestationDataFromNewPrysm(newPrysm *eth.AttestationData) *AttestationData {
	return &AttestationData{
		Slot:            uint64(newPrysm.Slot),
		CommitteeIndex:  uint64(newPrysm.CommitteeIndex),
		BeaconBlockRoot: newPrysm.BeaconBlockRoot,
		Source: &Checkpoint{
			Epoch: uint64(newPrysm.Source.Epoch),
			Root:  newPrysm.Source.Root,
		},
		Target: &Checkpoint{
			Epoch: uint64(newPrysm.Target.Epoch),
			Root:  newPrysm.Target.Root,
		},
	}
}

func (m *AttestationData) ToNewPrysm() *eth.AttestationData {
	return &eth.AttestationData{
		Slot:            types.Slot(m.Slot),
		CommitteeIndex:  types.CommitteeIndex(m.CommitteeIndex),
		BeaconBlockRoot: m.BeaconBlockRoot,
		Source: &eth.Checkpoint{
			Epoch: types.Epoch(m.Source.Epoch),
			Root:  m.Source.Root,
		},
		Target: &eth.Checkpoint{
			Epoch: types.Epoch(m.Target.Epoch),
			Root:  m.Target.Root,
		},
	}
}

func NewAggregationAndProofFromNewPrysm(newPrysm *eth.AggregateAttestationAndProof) *AggregateAttestationAndProof {
	return &AggregateAttestationAndProof{
		AggregatorIndex: uint64(newPrysm.AggregatorIndex),
		SelectionProof:  newPrysm.SelectionProof,
		Aggregate: &Attestation{
			AggregationBits: newPrysm.Aggregate.AggregationBits,
			Data:            NewAttestationDataFromNewPrysm(newPrysm.Aggregate.Data),
			Signature:       newPrysm.Aggregate.Signature,
		},
	}
}

func (m *AggregateAttestationAndProof) ToNewPrysm() *eth.AggregateAttestationAndProof {
	return &eth.AggregateAttestationAndProof{
		AggregatorIndex: types.ValidatorIndex(m.AggregatorIndex),
		SelectionProof:  m.SelectionProof,
		Aggregate: &eth.Attestation{
			AggregationBits: m.Aggregate.AggregationBits,
			Data:            m.Aggregate.Data.ToNewPrysm(),
			Signature:       m.Aggregate.Signature,
		},
	}
}
