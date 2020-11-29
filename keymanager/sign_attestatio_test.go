package keymanager_test

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/bloxapp/key-vault/keymanager"
	"github.com/hashicorp/vault/sdk/logical"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSignAttestation(t *testing.T) {
	domain := make([]byte, 32)
	rand.Read(domain)
	data := &ethpb.AttestationData{
		Slot:            10,
		CommitteeIndex:  10,
		BeaconBlockRoot: []byte{1, 2, 3},
		Source: &ethpb.Checkpoint{
			Epoch: 101010,
			Root:  []byte{2, 3, 4},
		},
		Target: &ethpb.Checkpoint{
			Epoch: 202020,
			Root:  []byte{5, 6, 7},
		},
	}

	beaconState, privKeys := testutil.DeterministicGenesisState(t, 100)
	block, err := testutil.GenerateFullBlock(beaconState, privKeys, nil, 0)
	require.NoError(t, err)

	expectedSignature, err := bls.SignatureFromBytes(block.GetSignature())
	require.NoError(t, err)

	actualSignature := hex.EncodeToString(block.GetSignature())

	var protect sync.Mutex
	var currentMethod http.HandlerFunc
	s := newTestRemoteWallet(func(writer http.ResponseWriter, request *http.Request) {
		currentMethod(writer, request)
	})
	defer s.Close()

	wallet, err := keymanager.NewKeyManager(logrus.NewEntry(logrus.New()), &keymanager.Config{
		Location:    s.URL,
		AccessToken: DefaultAccessToken,
		PubKey:      DefaultAccountPublicKey,
		Network:     "pyrmont",
	})
	require.NoError(t, err)

	runTest := func(t *testing.T, statusCode int, signature []byte, f func(wallet *keymanager.KeyManager)) {
		protect.Lock()
		currentMethod = func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, "/v1/ethereum/pyrmont/accounts/sign-attestation", request.URL.Path)

			var req keymanager.SignAttestationRequest
			require.NoError(t, json.NewDecoder(request.Body).Decode(&req))

			require.Equal(t, DefaultAccountPublicKey, req.PubKey)
			require.Equal(t, hex.EncodeToString(domain), req.Domain)
			require.Equal(t, int(data.GetSlot()), int(req.Slot))
			require.Equal(t, int(data.GetCommitteeIndex()), int(req.CommitteeIndex))
			require.Equal(t, hex.EncodeToString(data.GetBeaconBlockRoot()), req.BeaconBlockRoot)
			require.Equal(t, int(data.GetSource().GetEpoch()), int(req.SourceEpoch))
			require.Equal(t, hex.EncodeToString(data.GetSource().GetRoot()), req.SourceRoot)
			require.Equal(t, int(data.GetTarget().GetEpoch()), int(req.TargetEpoch))
			require.Equal(t, hex.EncodeToString(data.GetTarget().GetRoot()), req.TargetRoot)

			if statusCode == http.StatusOK {
				respBody := &logical.Response{
					Data: map[string]interface{}{
						"signature": string(signature),
					},
				}
				require.NoError(t, json.NewEncoder(writer).Encode(respBody))
			} else {
				writer.WriteHeader(statusCode)
			}
		}
		protect.Unlock()

		f(wallet)
	}

	t.Run("successfully signed data", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(domain, data)
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.Equal(t, expectedSignature, actualSignature)
		})
	})

	t.Run("rejects with denied", func(t *testing.T) {
		runTest(t, http.StatusUnauthorized, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(domain, data)
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with failed", func(t *testing.T) {
		runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(domain, data)
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with invalid signature", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte("invalid"), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(domain, data)
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})
}
