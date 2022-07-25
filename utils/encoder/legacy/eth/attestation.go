package eth

import (
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// NewAttestationDataFromNewPrysm new AttestationData struct
func NewAttestationDataFromNewPrysm(newPrysm *eth.AttestationData) *AttestationData {
	ret := &AttestationData{
		Slot:            uint64(newPrysm.Slot),
		CommitteeIndex:  uint64(newPrysm.CommitteeIndex),
		BeaconBlockRoot: newPrysm.BeaconBlockRoot,
	}

	if newPrysm.Source != nil {
		ret.Source = &Checkpoint{
			Epoch: uint64(newPrysm.Source.Epoch),
			Root:  newPrysm.Source.Root,
		}
	}
	if newPrysm.Target != nil {
		ret.Target = &Checkpoint{
			Epoch: uint64(newPrysm.Target.Epoch),
			Root:  newPrysm.Target.Root,
		}
	}
	return ret
}

// ToNewPrysm returns new AttestationData struct
func (m *AttestationData) ToNewPrysm() *eth.AttestationData {
	ret := &eth.AttestationData{}
	m.ToPrysm(ret)
	return ret
}

// ToPrysm returns AttestationData struct
func (m *AttestationData) ToPrysm(ret *eth.AttestationData) {
	ret.Slot = types.Slot(m.Slot)
	ret.CommitteeIndex = types.CommitteeIndex(m.CommitteeIndex)
	ret.BeaconBlockRoot = m.BeaconBlockRoot
	if m.Source != nil {
		ret.Source = &eth.Checkpoint{
			Epoch: types.Epoch(m.Source.Epoch),
			Root:  m.Source.Root,
		}
	}
	if m.Target != nil {
		ret.Target = &eth.Checkpoint{
			Epoch: types.Epoch(m.Target.Epoch),
			Root:  m.Target.Root,
		}
	}
}

// NewAggregationAndProofFromNewPrysm returns new AggregateAttestationAndProof  struct
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

// ToNewPrysm returns AggregateAttestationAndProof
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
