package encoderv2

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/prysmaticlabs/go-bitfield"

	types "github.com/prysmaticlabs/eth2-types"

	"github.com/stretchr/testify/require"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestV2(t *testing.T) {
	t.Run("attestation data", func(t *testing.T) {
		attestationDataByts := _byteArray("000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b")
		attData := &eth.AttestationData{}
		require.NoError(t, attData.UnmarshalSSZ(attestationDataByts))

		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestAttestationData{AttestationData: attData},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		byts, err = decoded.GetAttestationData().MarshalSSZ()
		require.NoError(t, err)
		require.EqualValues(t, attestationDataByts, byts)
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("beacon block", func(t *testing.T) {
		dataByts := _byteArray("010000000000000055000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776badd5cb7e6a4bffd8ce7fe9697aed511661861e312ad546dcf5480159698f47a554000000a2c156a4bc9439f1d85f922f2abaa96e830f1c526101211bdb7d16f4ad9490a0302fc5adb089c05b5f16fd465962f47c04fc2b81a94d135a07c1613db61511c17284b51fafab984e56d3411e16e45f5068f146d9412f91d31ab0f237eac3d745a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a8000000000000000a4e544482366bc9d5386f1cd0c4bf837327605620bf40c5514d51dfcadd14a4a0000000000000000000000000000000000000000000000000000000000000000dc000000dc000000dc000000c5010000c501000004000000e4000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b97b6f271ac364b041cd465f32fa7ffa19f5a811f1e6e14713f93e06537ef827d382bac72f0990b84f83cd9bbe0062815020086bf27b9ced172cc6add8ba5197991cf634d18666f5d43df6f09180ce20a357e4d05b2784409e32147f1042986e31f")
		data := &eth.BeaconBlock{}
		require.NoError(t, data.UnmarshalSSZ(dataByts))

		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestBlock{Block: data},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		byts, err = decoded.GetBlock().MarshalSSZ()
		require.NoError(t, err)
		require.EqualValues(t, dataByts, byts)
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("beacon block altair", func(t *testing.T) {
		dataByts := _byteArray("01000000000000001c00000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c83387dd0abb441a3c16886c8144098cb4cac5e363516f329c368550094fd7ff754000000b1e2f27dfac80e4f1bce84adf11acf6cdbb0d8e59a575c9795020e614eb3aa29634108c0559c04ce02b93fc9a5a8daf60485ebac039864c79d51bef54915aa8c45cbcde3215f14962be196a6b8648851c35b4a804ce8d5fb6c5ff49800ef7740685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb732000000000000000685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb7300000000000000000000000000000000000000000000000000000000000000007c0100007c0100007c0100006502000065020000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa2ec291dd5e91096ae48b3659a7ac59567a48c030bb6ac9435d6d44ef39f3f664742f35b38cd6e41ade9ed417183cc0c0b407dfea8627ccc2275fc82ab3d2182e58a037eb144811d741d18894698396efde2b7873c2db9b712e03dfcd03705ef04000000e400000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669cb62ce3f28e8731dce73d5761fdc5e30383d42a022d6e939974d0586d82270f79b38b86d17237e4241a761e239c594e7a0d4ef731470001be3b125ba515f8f215f9309a9ba12653bf9d704a4125865b9775c8a65223e3ca027781175200a2d24403")
		data := &eth.BeaconBlockAltair{}
		require.NoError(t, data.UnmarshalSSZ(dataByts))

		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestBlockV2{BlockV2: data},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		byts, err = decoded.GetBlockV2().MarshalSSZ()
		require.NoError(t, err)
		require.EqualValues(t, dataByts, byts)
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("attestation aggregation", func(t *testing.T) {
		dataByts := _byteArray("01000000000000001c00000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c83387dd0abb441a3c16886c8144098cb4cac5e363516f329c368550094fd7ff754000000b1e2f27dfac80e4f1bce84adf11acf6cdbb0d8e59a575c9795020e614eb3aa29634108c0559c04ce02b93fc9a5a8daf60485ebac039864c79d51bef54915aa8c45cbcde3215f14962be196a6b8648851c35b4a804ce8d5fb6c5ff49800ef7740685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb732000000000000000685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb7300000000000000000000000000000000000000000000000000000000000000007c0100007c0100007c0100006502000065020000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa2ec291dd5e91096ae48b3659a7ac59567a48c030bb6ac9435d6d44ef39f3f664742f35b38cd6e41ade9ed417183cc0c0b407dfea8627ccc2275fc82ab3d2182e58a037eb144811d741d18894698396efde2b7873c2db9b712e03dfcd03705ef04000000e400000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669cb62ce3f28e8731dce73d5761fdc5e30383d42a022d6e939974d0586d82270f79b38b86d17237e4241a761e239c594e7a0d4ef731470001be3b125ba515f8f215f9309a9ba12653bf9d704a4125865b9775c8a65223e3ca027781175200a2d24403")
		data := &eth.BeaconBlockAltair{}
		require.NoError(t, data.UnmarshalSSZ(dataByts))

		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestBlockV2{BlockV2: data},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		byts, err = decoded.GetBlockV2().MarshalSSZ()
		require.NoError(t, err)
		require.EqualValues(t, dataByts, byts)
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("slot", func(t *testing.T) {
		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestSlot{Slot: 2},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		require.EqualValues(t, types.Slot(2), decoded.GetSlot())
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("epoch", func(t *testing.T) {
		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestEpoch{Epoch: 2},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		require.EqualValues(t, types.Epoch(2), decoded.GetEpoch())
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.PublicKey)
	})
	t.Run("sync committee", func(t *testing.T) {
		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object:          &models.SignRequestSyncCommitteeMessage{Root: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.GetSyncCommitteeMessage())
	})
	t.Run("sync aggregator", func(t *testing.T) {
		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object: &models.SignRequestSyncAggregatorSelectionData{
				SyncAggregatorSelectionData: &eth.SyncAggregatorSelectionData{
					Slot:              types.Slot(12),
					SubcommitteeIndex: 44,
				},
			},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		require.EqualValues(t, 12, decoded.GetSyncAggregatorSelectionData().Slot)
		require.EqualValues(t, 44, decoded.GetSyncAggregatorSelectionData().SubcommitteeIndex)
	})
	t.Run("sync contribution proof", func(t *testing.T) {
		req := &models.SignRequest{
			PublicKey:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
			SigningRoot:     make([]byte, 32),
			SignatureDomain: make([]byte, 32),
			Object: &models.SignRequestContributionAndProof{
				ContributionAndProof: &eth.ContributionAndProof{
					AggregatorIndex: types.ValidatorIndex(12),
					SelectionProof:  make([]byte, 96),
					Contribution: &eth.SyncCommitteeContribution{
						Slot:              11,
						BlockRoot:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
						SubcommitteeIndex: 12,
						AggregationBits:   bitfield.NewBitvector128(),
						Signature:         make([]byte, 96),
					},
				},
			},
		}

		enc := New()
		byts, err := enc.Encode(req)
		require.NoError(t, err)

		decoded := &models.SignRequest{}
		require.NoError(t, enc.Decode(byts, decoded))
		require.EqualValues(t, 12, decoded.GetContributionAndProof().AggregatorIndex)
		require.EqualValues(t, make([]byte, 96), decoded.GetContributionAndProof().SelectionProof)
		require.EqualValues(t, 11, decoded.GetContributionAndProof().Contribution.Slot)
		require.EqualValues(t, 12, decoded.GetContributionAndProof().Contribution.SubcommitteeIndex)
		require.EqualValues(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1}, decoded.GetContributionAndProof().Contribution.BlockRoot)
	})
}
