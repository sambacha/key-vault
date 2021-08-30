package encoder

import (
	"encoding/hex"
	"testing"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestLegacyAttestationDataMarshal(t *testing.T) {
	sszAttData := _byteArray("e7030000000000000000000000000000626c6f636b526f6f74000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000000000000000000000000")
	originalSSZ := eth.AttestationData{}
	require.NoError(t, originalSSZ.UnmarshalSSZ(sszAttData))

	marshalByts := _byteArray("08e7071a20626c6f636b526f6f74000000000000000000000000000000000000000000000022240864122000000000000000000000000000000000000000000000000000000000000000002a2508c80112200000000000000000000000000000000000000000000000000000000000000000")
	marshaled := &eth.AttestationData{}
	require.NoError(t, LegacyAttestationDataUnMarshal(marshaled, marshalByts))

	// marshal and unmarshal
	remarshalByts, err := LegacyAttestationDataMarshal(marshaled)
	require.NoError(t, err)
	remarshaled := &eth.AttestationData{}
	require.NoError(t, LegacyAttestationDataUnMarshal(remarshaled, remarshalByts))

	// verify
	require.EqualValues(t, originalSSZ.Slot, marshaled.Slot)
	require.EqualValues(t, remarshaled.Slot, marshaled.Slot)

	require.EqualValues(t, originalSSZ.Target.Epoch, marshaled.Target.Epoch)
	require.EqualValues(t, remarshaled.Target.Epoch, marshaled.Target.Epoch)

	require.EqualValues(t, originalSSZ.Target.Root, marshaled.Target.Root)
	require.EqualValues(t, remarshaled.Target.Root, marshaled.Target.Root)

	require.EqualValues(t, originalSSZ.Source.Epoch, marshaled.Source.Epoch)
	require.EqualValues(t, remarshaled.Source.Epoch, marshaled.Source.Epoch)

	require.EqualValues(t, originalSSZ.Source.Root, marshaled.Source.Root)
	require.EqualValues(t, remarshaled.Source.Root, marshaled.Source.Root)

	require.EqualValues(t, originalSSZ.BeaconBlockRoot, marshaled.BeaconBlockRoot)
	require.EqualValues(t, remarshaled.BeaconBlockRoot, marshaled.BeaconBlockRoot)
}
