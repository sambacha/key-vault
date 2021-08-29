package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

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
		AggregatorIndex: 0,
		Aggregate: &ethpb.Attestation{
			Data: &ethpb.AttestationData{
				BeaconBlockRoot: make([]byte, 32),
				Target:          &ethpb.Checkpoint{Root: make([]byte, 32)},
				Source:          &ethpb.Checkpoint{Root: make([]byte, 32)},
			},
			Signature:       make([]byte, 96),
			AggregationBits: make([]byte, 1),
		},
		SelectionProof: make([]byte, 96),
	}

	req := &validatorpb.SignRequest{
		PublicKey:       _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
		SigningRoot:     nil,
		SignatureDomain: _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
		Object:          &validatorpb.SignRequest_AggregateAttestationAndProof{AggregateAttestationAndProof: agg},
	}

	if undefinedPubKey {
		req.PublicKey = _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcd")
	}

	byts, _ := json.Marshal(req)
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}
}

func TestSignAttestation(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Sign Attestation in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	})

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

	t.Run("Successfully Sign Proposal", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		data := basicProposalData()
		req.Data = data
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	t.Run("Sign Proposal in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	})

	t.Run("Sign Proposal of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicProposalDataWithOps(true, false, false, false)
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

	t.Run("Sign Aggregation in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign: failed to open key vault: wallet not found")
	})

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
