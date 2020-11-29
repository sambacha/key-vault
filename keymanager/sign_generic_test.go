package keymanager_test

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/bloxapp/key-vault/keymanager"
	"github.com/bloxapp/key-vault/utils/bytex"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSignGeneric(t *testing.T) {
	domain := make([]byte, 32)
	rand.Read(domain)

	data := make([]byte, 32)
	rand.Read(data)

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
			require.Equal(t, "/v1/ethereum/pyrmont/accounts/sign-aggregation", request.URL.Path)

			var req keymanager.SignAggregationRequest
			require.NoError(t, json.NewDecoder(request.Body).Decode(&req))

			require.Equal(t, DefaultAccountPublicKey, req.PubKey)
			require.Equal(t, hex.EncodeToString(domain), req.Domain)
			require.Equal(t, hex.EncodeToString(data), req.DataToSign)

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
			actualSignature, err := wallet.SignGeneric(data, bytex.ToBytes32(domain))
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.Equal(t, expectedSignature, actualSignature)
		})
	})

	t.Run("rejects with denied", func(t *testing.T) {
		runTest(t, http.StatusUnauthorized, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(data, bytex.ToBytes32(domain))
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with failed", func(t *testing.T) {
		runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(data, bytex.ToBytes32(domain))
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with invalid signature", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte("invalid"), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(data, bytex.ToBytes32(domain))
			require.True(t, keymanager.IsGenericError(err))
			require.Nil(t, actualSignature)
		})
	})
}
