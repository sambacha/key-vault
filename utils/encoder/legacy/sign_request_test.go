package legacy

import (
	"testing"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	"github.com/stretchr/testify/require"
)

func TestSignRequestJsonMarshal(t *testing.T) {
	t.Run("with attestation data", func(t *testing.T) {
		marshalByts := _byteArray("7b227075626c69635f6b6579223a22736b585750542b646a714741656d4b66797873796a4c545651764e6150567645654c344e38346e6433584576784d6757756a2f7436616c6a494b35724a4b6659222c227369676e696e675f726f6f74223a22414141414145553172537a4935745435544e3648423674652b6166535034684e2b45776c4d594d714b31343d222c227369676e61747572655f646f6d61696e223a22414141414145553172537a4935745435544e3648423674652b6166535034684e2b45776c4d594d714b31343d222c224f626a656374223a7b226174746573746174696f6e5f64617461223a7b22736c6f74223a3939392c22626561636f6e5f626c6f636b5f726f6f74223a22596d78765932745362323930414141414141414141414141414141414141414141414141414141414141413d222c22736f75726365223a7b2265706f6368223a3130302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a3230302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d7d7d7d")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestJsonUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestJsonMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("626c6f636b526f6f740000000000000000000000000000000000000000000000s"), marshaled.GetAttestationData().BeaconBlockRoot)
	})
	t.Run("with beacon block", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5eaa06ed03080c100a1a20c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf22220dad936bf55f3c4743a6659eb4e9b50f791aa781a9b470b377e156b358e27b2052aa2030a608d00ac04a601da595c750b2177f6941785937a014c3783b6eb5663992f591361e73abd2f18f3401679cbfa2ac7dcb5460464eed6368507135e73704336ecadf27dc0f092619754b355a7b796e88e89c76a117ccfcf71e2a5d42cf02429acb04012460a20685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb7310201a20685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb731a20000000000000000000000000000000000000000000000000000000000000000032d3010a0103126c080b1a20c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf22222122000000000000000000000000000000000000000000000000000000000000000002a221220c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf21a60aff989b94be8ed7d2b49f493328c6e4957c3a0b32568bc7f1a68dfb28f8c37c7254c8817281ed4770ea4b555f765e25b0887b515c8eb48e56ae94fb51ab005b7d60639ea96203c03233fe3fbdd45b4791bab4f2a8e0012ae3b43a4455849e533")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("dad936bf55f3c4743a6659eb4e9b50f791aa781a9b470b377e156b358e27b205s"), marshaled.GetBlock().StateRoot)
	})
	t.Run("with aggregation", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5eba06dd011ada010a020010127208e7071a20626c6f636b526f6f74000000000000000000000000000000000000000000000022240864122000000000000000000000000000000000000000000000000000000000000000002a2508c801122000000000000000000000000000000000000000000000000000000000000000001a60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("626c6f636b526f6f740000000000000000000000000000000000000000000000s"), marshaled.GetAggregateAttestationAndProof().Aggregate.Data.BeaconBlockRoot)
	})
	t.Run("with epoch", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5ed0060c")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, 12, marshaled.GetEpoch())
	})
	t.Run("with slot", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5ec8060c")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, 12, marshaled.GetSlot())
	})
}

func TestSignRequestMarshal(t *testing.T) {
	t.Run("with attestation data", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5eb2067208e7071a20626c6f636b526f6f74000000000000000000000000000000000000000000000022240864122000000000000000000000000000000000000000000000000000000000000000002a2508c80112200000000000000000000000000000000000000000000000000000000000000000")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("626c6f636b526f6f740000000000000000000000000000000000000000000000s"), marshaled.GetAttestationData().BeaconBlockRoot)
	})
	t.Run("with beacon block", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5eaa06ed03080c100a1a20c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf22220dad936bf55f3c4743a6659eb4e9b50f791aa781a9b470b377e156b358e27b2052aa2030a608d00ac04a601da595c750b2177f6941785937a014c3783b6eb5663992f591361e73abd2f18f3401679cbfa2ac7dcb5460464eed6368507135e73704336ecadf27dc0f092619754b355a7b796e88e89c76a117ccfcf71e2a5d42cf02429acb04012460a20685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb7310201a20685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb731a20000000000000000000000000000000000000000000000000000000000000000032d3010a0103126c080b1a20c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf22222122000000000000000000000000000000000000000000000000000000000000000002a221220c26a0752c187df9e835bcea4b2cdf494360860ff589e7c7f30be21917f111bf21a60aff989b94be8ed7d2b49f493328c6e4957c3a0b32568bc7f1a68dfb28f8c37c7254c8817281ed4770ea4b555f765e25b0887b515c8eb48e56ae94fb51ab005b7d60639ea96203c03233fe3fbdd45b4791bab4f2a8e0012ae3b43a4455849e533")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("dad936bf55f3c4743a6659eb4e9b50f791aa781a9b470b377e156b358e27b205s"), marshaled.GetBlock().StateRoot)
	})
	t.Run("with aggregation", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5eba06dd011ada010a020010127208e7071a20626c6f636b526f6f74000000000000000000000000000000000000000000000022240864122000000000000000000000000000000000000000000000000000000000000000002a2508c801122000000000000000000000000000000000000000000000000000000000000000001a60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, _byteArray("626c6f636b526f6f740000000000000000000000000000000000000000000000s"), marshaled.GetAggregateAttestationAndProof().Aggregate.Data.BeaconBlockRoot)
	})
	t.Run("with epoch", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5ed0060c")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, 12, marshaled.GetEpoch())
	})
	t.Run("with slot", func(t *testing.T) {
		marshalByts := _byteArray("0a30b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d81220000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e1a20000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5ec8060c")
		marshaled := &validatorpb.SignRequest{}
		require.NoError(t, LegacySignRequestUnMarshal(marshaled, marshalByts))

		// marshal and unmarshal
		remarshalByts, err := LegacySignRequestMarshal(marshaled)
		require.NoError(t, err)
		require.EqualValues(t, marshalByts, remarshalByts)

		// verify
		require.EqualValues(t, _byteArray("000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"), marshaled.SigningRoot)
		require.EqualValues(t, 12, marshaled.GetSlot())
	})
}
