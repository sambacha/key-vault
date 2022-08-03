package backend

import (
	"bytes"
	"context"
	"encoding/hex"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/prysmaticlabs/go-bitfield"

	types "github.com/prysmaticlabs/eth2-types"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func setupStorageWithWalletAndAccounts(storage logical.Storage) error {
	_, err := baseHashicorpStorage(context.Background(), storage)
	return err
}

func basicAggregationAndProofData() map[string]interface{} {
	return basicAggregationAndProofDataWithOps(false)
}

func basicAggregationAndProofDataWithOps(undefinedPubKey bool) map[string]interface{} {
	agg := &ethpb.AggregateAttestationAndProof{
		AggregatorIndex: types.ValidatorIndex(1),
		SelectionProof:  make([]byte, 96),
		Aggregate: &ethpb.Attestation{
			AggregationBits: bitfield.NewBitlist(12),
			Signature:       make([]byte, 96),
			Data: &ethpb.AttestationData{
				Slot:            types.Slot(1),
				CommitteeIndex:  types.CommitteeIndex(12),
				BeaconBlockRoot: make([]byte, 32),
				Source: &ethpb.Checkpoint{
					Epoch: types.Epoch(1),
					Root:  make([]byte, 32),
				},
				Target: &ethpb.Checkpoint{
					Epoch: types.Epoch(1),
					Root:  make([]byte, 32),
				},
			},
		},
	}

	req := &models.SignRequest{
		PublicKey:       _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
		SigningRoot:     nil,
		SignatureDomain: _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
		Object:          &models.SignRequestAggregateAttestationAndProof{AggregateAttestationAndProof: agg},
	}

	if undefinedPubKey {
		req.PublicKey = _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcd")
	}

	byts, _ := encoderv2.New().Encode(req)
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}
}

func TestSignAttestation(t *testing.T) {
	b, _ := getBackend(t)

	//t.Run("Sign Attestation in non existing key vault", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
	//	setupBaseStorage(t, req)
	//	_, err := b.HandleRequest(context.Background(), req)
	//	require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	//})

	t.Run("Sign Attestation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicAttestationDataWithOps(true, false, false, false, false)
		resp, err := b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: account not found")
		require.Nil(t, resp)
	})
}

func TestSignProposal(t *testing.T) {
	b, _ := getBackend(t)

	withEachBlockVersion(t, "Successfully Sign Proposal", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		data := basicProposalData(blockVersion)
		req.Data = data
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	//t.Run("Sign Proposal in non existing key vault", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
	//	setupBaseStorage(t, req)
	//	_, err := b.HandleRequest(context.Background(), req)
	//	require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	//})

	withEachBlockVersion(t, "Sign Proposal of unknown account", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicProposalDataWithOps(blockVersion, true, false, false, false)
		resp, err := b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: account not found")
		require.Nil(t, resp)
	})
}

func TestSignAggregation(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Aggregation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicAggregationAndProofData()
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
	})

	//t.Run("Sign Aggregation in non existing key vault", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
	//	setupBaseStorage(t, req)
	//	_, err := b.HandleRequest(context.Background(), req)
	//	require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	//})

	t.Run("Sign Aggregation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicAggregationAndProofDataWithOps(true)
		resp, err := b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: account not found")
		require.Nil(t, resp)
	})
}

func TestValidateRequestedFeeRecipient(t *testing.T) {
	recipient := func(lastByte byte) string { return hexutil.Encode(append(bytes.Repeat([]byte{0}, 19), lastByte)) }
	pubKey := func(lastByte byte) string { return hexutil.Encode(append(bytes.Repeat([]byte{0}, 95), lastByte)) }

	tests := []struct {
		name               string
		configured         FeeRecipients
		requestedValidator string
		requestedRecipient string
		expectedErr        error
	}{
		{
			name: "Good",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(1),
			requestedRecipient: recipient(1),
			expectedErr:        nil,
		},
		{
			name: "Default",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(3),
			requestedRecipient: recipient(3),
			expectedErr:        nil,
		},
		{
			name: "Wrong default",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(3),
			requestedRecipient: recipient(4),
			expectedErr:        ErrFeeRecipientDiffers,
		},
		{
			name: "Differs from configuration",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(1),
			requestedRecipient: recipient(2),
			expectedErr:        ErrFeeRecipientDiffers,
		},
		{
			name: "No such validator",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
			},
			requestedValidator: pubKey(3),
			requestedRecipient: recipient(2),
			expectedErr:        ErrFeeRecipientNotSet,
		},
		{
			name: "No such recipient",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(1),
			requestedRecipient: recipient(4),
			expectedErr:        ErrFeeRecipientDiffers,
		},
		{
			name: "Other validator's recipient",
			configured: FeeRecipients{
				pubKey(1): recipient(1),
				pubKey(2): recipient(2),
				"default": recipient(3),
			},
			requestedValidator: pubKey(1),
			requestedRecipient: recipient(2),
			expectedErr:        ErrFeeRecipientDiffers,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateRequestedFeeRecipient(hexutil.MustDecode(test.requestedValidator), test.configured, hexutil.MustDecode(test.requestedRecipient))
			require.Equal(t, test.expectedErr, err)
		})
	}
}
