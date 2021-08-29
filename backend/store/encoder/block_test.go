package encoder

import (
	"testing"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

/**
blk := &ethpb.BeaconBlock{
		Slot:12,
		ProposerIndex:122,
		ParentRoot: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
		StateRoot:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
		Body: &ethpb.BeaconBlockBody{
			RandaoReveal:         make([]byte, 96),
			Eth1Data:             &ethpb.Eth1Data{
				DepositRoot:          []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
				DepositCount:         12,
				BlockHash:            []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
			},
			Graffiti:             []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
			ProposerSlashings:    []*ethpb.ProposerSlashing{
				{
					Header_1:             &ethpb.SignedBeaconBlockHeader{
						Header:              &ethpb.BeaconBlockHeader{
							Slot:12,
							ProposerIndex:33,
							ParentRoot: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							StateRoot: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							BodyRoot:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						},
						Signature:            make([]byte, 96),
					},
					Header_2:               &ethpb.SignedBeaconBlockHeader{
						Header:              &ethpb.BeaconBlockHeader{
							Slot:12,
							ProposerIndex:33,
							ParentRoot: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							StateRoot: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							BodyRoot:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						},
						Signature:            make([]byte, 96),
					},
				},
			},
			AttesterSlashings:    []*ethpb.AttesterSlashing{
				{
					Attestation_1: &ethpb.IndexedAttestation{
						AttestingIndices:     []uint64{1,2,3},
						Data:                 &ethpb.AttestationData{
							Slot:12,
							CommitteeIndex:1,
							BeaconBlockRoot:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							Source: &ethpb.Checkpoint{
								Epoch:1,
								Root:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							},
							Target:&ethpb.Checkpoint{
								Epoch:1,
								Root:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
							},
						},
						Signature:            make([]byte, 96),
					},
				},
			},
			Attestations:         []*ethpb.Attestation{
				{
					AggregationBits: bitfield.NewBitlist(12),
					Data:&ethpb.AttestationData{
						Slot:12,
						CommitteeIndex:1,
						BeaconBlockRoot:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						Source: &ethpb.Checkpoint{
							Epoch:1,
							Root:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						},
						Target:&ethpb.Checkpoint{
							Epoch:1,
							Root:[]byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						},
					},
					Signature: make([]byte, 96),
				},
			},
			Deposits:             []*ethpb.Deposit{
				{
					Proof: make([][]byte,33),
					Data: &ethpb.Deposit_Data{
						PublicKey:             make([]byte, 48),
						WithdrawalCredentials: []byte{1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,1,2,3,4,5,9,9},
						Amount:                32,
						Signature:             make([]byte, 96),
					},
				},
			},
			VoluntaryExits:       []*ethpb.SignedVoluntaryExit{
				{
						Exit: &ethpb.VoluntaryExit{
							Epoch: 12,
							ValidatorIndex:33,
						},
						Signature: make([]byte, 96),
				},
			},
		},
	}
*/
func TestLegacyBlockMarshal(t *testing.T) {
	sszAttData := _byteArray("e7030000000000000000000000000000626c6f636b526f6f74000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000000000000000000000000")
	originalSSZ := eth.AttestationData{}
	require.NoError(t, originalSSZ.UnmarshalSSZ(sszAttData))

	marshalByts := _byteArray("08e7071a20626c6f636b526f6f74000000000000000000000000000000000000000000000022240864122000000000000000000000000000000000000000000000000000000000000000002a2508c80112200000000000000000000000000000000000000000000000000000000000000000")
	marshaled := &eth.AttestationData{}
	require.NoError(t, LegacyAttestationDataUnMarshal(marshaled, marshalByts))

	// verify
	require.EqualValues(t, originalSSZ.Slot, marshaled.Slot)
	require.EqualValues(t, originalSSZ.Target.Epoch, marshaled.Target.Epoch)
	require.EqualValues(t, originalSSZ.Target.Root, marshaled.Target.Root)
	require.EqualValues(t, originalSSZ.Source.Epoch, marshaled.Source.Epoch)
	require.EqualValues(t, originalSSZ.Source.Root, marshaled.Source.Root)
	require.EqualValues(t, originalSSZ.BeaconBlockRoot, marshaled.BeaconBlockRoot)
}
