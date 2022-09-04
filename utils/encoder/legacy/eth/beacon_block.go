package eth

import (
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// ToNewPrysm returns new BeaconBlock
func (m *BeaconBlock) ToNewPrysm() *eth.BeaconBlock {
	ret := &eth.BeaconBlock{}
	m.ToPrysm(ret)
	return ret
}

// ToPrysm returns new BeaconBlock
func (m *BeaconBlock) ToPrysm(ret *eth.BeaconBlock) {
	ret.Slot = types.Slot(m.Slot)
	ret.ProposerIndex = types.ValidatorIndex(m.ProposerIndex)
	ret.ParentRoot = m.ParentRoot
	ret.StateRoot = m.StateRoot

	if m.Body == nil {
		return
	}

	ret.Body = &eth.BeaconBlockBody{
		RandaoReveal: m.Body.RandaoReveal,
		Eth1Data: &eth.Eth1Data{
			DepositRoot:  m.Body.Eth1Data.DepositRoot,
			DepositCount: m.Body.Eth1Data.DepositCount,
			BlockHash:    m.Body.Eth1Data.BlockHash,
		},
		Graffiti:          m.Body.Graffiti,
		ProposerSlashings: make([]*eth.ProposerSlashing, 0),
		AttesterSlashings: make([]*eth.AttesterSlashing, 0),
		Attestations:      make([]*eth.Attestation, 0),
		Deposits:          make([]*eth.Deposit, 0),
		VoluntaryExits:    make([]*eth.SignedVoluntaryExit, 0),
	}

	for _, prop := range m.Body.ProposerSlashings {
		ret.Body.ProposerSlashings = append(ret.Body.ProposerSlashings, &eth.ProposerSlashing{
			Header_1: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          types.Slot(prop.Header_1.Header.Slot),
					ProposerIndex: types.ValidatorIndex(prop.Header_1.Header.ProposerIndex),
					ParentRoot:    prop.Header_1.Header.ParentRoot,
					StateRoot:     prop.Header_1.Header.StateRoot,
					BodyRoot:      prop.Header_1.Header.BodyRoot,
				},
				Signature: prop.Header_1.Signature,
			},
			Header_2: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          types.Slot(prop.Header_2.Header.Slot),
					ProposerIndex: types.ValidatorIndex(prop.Header_2.Header.ProposerIndex),
					ParentRoot:    prop.Header_2.Header.ParentRoot,
					StateRoot:     prop.Header_2.Header.StateRoot,
					BodyRoot:      prop.Header_2.Header.BodyRoot,
				},
				Signature: prop.Header_2.Signature,
			},
		})
	}

	for _, att := range m.Body.AttesterSlashings {
		ret.Body.AttesterSlashings = append(ret.Body.AttesterSlashings, &eth.AttesterSlashing{
			Attestation_1: &eth.IndexedAttestation{
				AttestingIndices: att.Attestation_1.AttestingIndices,
				Data:             att.Attestation_1.Data.ToNewPrysm(),
				Signature:        att.Attestation_1.Signature,
			},
			Attestation_2: &eth.IndexedAttestation{
				AttestingIndices: att.Attestation_2.AttestingIndices,
				Data:             att.Attestation_2.Data.ToNewPrysm(),
				Signature:        att.Attestation_2.Signature,
			},
		})
	}

	for _, att := range m.Body.Attestations {
		ret.Body.Attestations = append(ret.Body.Attestations, &eth.Attestation{
			AggregationBits: att.AggregationBits,
			Data:            att.Data.ToNewPrysm(),
			Signature:       att.Signature,
		})
	}

	for _, depo := range m.Body.Deposits {
		ret.Body.Deposits = append(ret.Body.Deposits, &eth.Deposit{
			Proof: depo.Proof,
			Data: &eth.Deposit_Data{
				PublicKey:             depo.Data.PublicKey,
				WithdrawalCredentials: depo.Data.WithdrawalCredentials,
				Amount:                depo.Data.Amount,
				Signature:             depo.Data.Signature,
			},
		})
	}

	for _, exit := range m.Body.VoluntaryExits {
		ret.Body.VoluntaryExits = append(ret.Body.VoluntaryExits, &eth.SignedVoluntaryExit{
			Exit: &eth.VoluntaryExit{
				Epoch:          types.Epoch(exit.Exit.Epoch),
				ValidatorIndex: types.ValidatorIndex(exit.Exit.ValidatorIndex),
			},
			Signature: exit.Signature,
		})
	}
}

// NewBeaconBlockFromNewPrysm returns BeaconBlock struct
func NewBeaconBlockFromNewPrysm(newPrysm *eth.BeaconBlock) *BeaconBlock {
	ret := &BeaconBlock{
		Slot:          uint64(newPrysm.Slot),
		ProposerIndex: uint64(newPrysm.ProposerIndex),
		ParentRoot:    newPrysm.ParentRoot,
		StateRoot:     newPrysm.StateRoot,
	}

	if newPrysm.Body == nil {
		return ret
	}

	ret.Body = &BeaconBlockBody{
		RandaoReveal: newPrysm.Body.RandaoReveal,
		Eth1Data: &Eth1Data{
			DepositRoot:  newPrysm.Body.Eth1Data.DepositRoot,
			DepositCount: newPrysm.Body.Eth1Data.DepositCount,
			BlockHash:    newPrysm.Body.Eth1Data.BlockHash,
		},
		Graffiti:          newPrysm.Body.Graffiti,
		ProposerSlashings: make([]*ProposerSlashing, 0),
		AttesterSlashings: make([]*AttesterSlashing, 0),
		Attestations:      make([]*Attestation, 0),
		Deposits:          make([]*Deposit, 0),
		VoluntaryExits:    make([]*SignedVoluntaryExit, 0),
	}

	for _, prop := range newPrysm.Body.ProposerSlashings {
		ret.Body.ProposerSlashings = append(ret.Body.ProposerSlashings, &ProposerSlashing{
			Header_1: &SignedBeaconBlockHeader{
				Header: &BeaconBlockHeader{
					Slot:          uint64(prop.Header_1.Header.Slot),
					ProposerIndex: uint64(prop.Header_1.Header.ProposerIndex),
					ParentRoot:    prop.Header_1.Header.ParentRoot,
					StateRoot:     prop.Header_1.Header.StateRoot,
					BodyRoot:      prop.Header_1.Header.BodyRoot,
				},
				Signature: prop.Header_1.Signature,
			},
			Header_2: &SignedBeaconBlockHeader{
				Header: &BeaconBlockHeader{
					Slot:          uint64(prop.Header_2.Header.Slot),
					ProposerIndex: uint64(prop.Header_2.Header.ProposerIndex),
					ParentRoot:    prop.Header_2.Header.ParentRoot,
					StateRoot:     prop.Header_2.Header.StateRoot,
					BodyRoot:      prop.Header_2.Header.BodyRoot,
				},
				Signature: prop.Header_2.Signature,
			},
		})
	}

	for _, att := range newPrysm.Body.AttesterSlashings {
		ret.Body.AttesterSlashings = append(ret.Body.AttesterSlashings, &AttesterSlashing{
			Attestation_1: &IndexedAttestation{
				AttestingIndices: att.Attestation_1.AttestingIndices,
				Data:             NewAttestationDataFromNewPrysm(att.Attestation_1.Data),
				Signature:        att.Attestation_1.Signature,
			},
			Attestation_2: &IndexedAttestation{
				AttestingIndices: att.Attestation_2.AttestingIndices,
				Data:             NewAttestationDataFromNewPrysm(att.Attestation_2.Data),
				Signature:        att.Attestation_2.Signature,
			},
		})
	}

	for _, att := range newPrysm.Body.Attestations {
		ret.Body.Attestations = append(ret.Body.Attestations, &Attestation{
			AggregationBits: att.AggregationBits,
			Data:            NewAttestationDataFromNewPrysm(att.Data),
			Signature:       att.Signature,
		})
	}

	for _, depo := range newPrysm.Body.Deposits {
		ret.Body.Deposits = append(ret.Body.Deposits, &Deposit{
			Proof: depo.Proof,
			Data: &Deposit_Data{
				PublicKey:             depo.Data.PublicKey,
				WithdrawalCredentials: depo.Data.WithdrawalCredentials,
				Amount:                depo.Data.Amount,
				Signature:             depo.Data.Signature,
			},
		})
	}

	for _, exit := range newPrysm.Body.VoluntaryExits {
		ret.Body.VoluntaryExits = append(ret.Body.VoluntaryExits, &SignedVoluntaryExit{
			Exit: &VoluntaryExit{
				Epoch:          uint64(exit.Exit.Epoch),
				ValidatorIndex: uint64(exit.Exit.ValidatorIndex),
			},
			Signature: exit.Signature,
		})
	}
	return ret
}
