package keymanager_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/keymanager"
)

func TestSignProposal(t *testing.T) {
	expectedSig := _byteArray("911ac2f6d74039279f16eee4cc46f4c6eea0ef9d18f0d9739b407c150c07ccb104c1c4b034ad46b25719bafc22fad05205975393000ea09636f5ce427814e2fe12ea72041099cc7f6ec249e504992dbf65e968ab448ddf4e124cbcbc722829b5")

	var protect sync.Mutex
	var currentMethod http.HandlerFunc
	s := newTestRemoteWallet(func(writer http.ResponseWriter, request *http.Request) {
		currentMethod(writer, request)
	})
	defer s.Close()

	km, err := keymanager.NewKeyManager(logrus.NewEntry(logrus.New()), &keymanager.Config{
		Location:    s.URL,
		AccessToken: DefaultAccessToken,
		PubKey:      "a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972",
		Network:     "prater",
	})
	require.NoError(t, err)

	runTest := func(t *testing.T, statusCode int, signature []byte, f func(wallet *keymanager.KeyManager)) {
		protect.Lock()
		currentMethod = func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, "/v1/ethereum/prater/accounts/sign", request.URL.Path)

			var reqBody map[string]interface{}
			require.NoError(t, json.NewDecoder(request.Body).Decode(&reqBody))

			// test un-marshaling the request
			val, ok := reqBody["sign_req"]
			require.True(t, ok)
			valByts, err := hex.DecodeString(val.(string))
			require.NoError(t, err)
			req := &v2.SignRequest{}
			require.NoError(t, req.Unmarshal(valByts))

			require.EqualValues(t, _byteArray("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972"), req.PublicKey)
			require.EqualValues(t, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), req.SignatureDomain)

			// root block
			actualRoot, err := req.Object.(*v2.SignRequest_Block).Block.HashTreeRoot()
			require.NoError(t, err)
			expectedRoot := _byteArray("c842b15864adc006334b251e1b692f21ae453a1524e8d2224d602a2822949190")
			require.EqualValues(t, expectedRoot, actualRoot[:])

			if statusCode == http.StatusOK {
				respBody := &logical.Response{
					Data: map[string]interface{}{
						"signature": hex.EncodeToString(signature),
					},
				}
				require.NoError(t, json.NewEncoder(writer).Encode(respBody))
			} else {
				writer.WriteHeader(statusCode)
			}
		}
		protect.Unlock()

		f(km)
	}

	t.Run("successfully signed data", func(t *testing.T) {
		runTest(t, http.StatusOK, expectedSig, func(km *keymanager.KeyManager) {
			actualSignature, err := km.Sign(context.Background(), testRequest(t))
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.EqualValues(t, expectedSig, actualSignature.Marshal())
		})
	})

	t.Run("provide wrong pubKey", func(t *testing.T) {
		runTest(t, http.StatusOK, expectedSig, func(km *keymanager.KeyManager) {
			req := testRequest(t)
			req.PublicKey = _byteArray("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186971")
			actualSignature, err := km.Sign(context.Background(), req)
			require.NotNil(t, err)
			require.True(t, keymanager.IsGenericError(err))
			require.EqualError(t, err, "{\"error\":\"no such key\"}")
			require.Nil(t, actualSignature)
		})
	})

	//t.Run("rejects with denied", func(t *testing.T) {
	//	runTest(t, http.StatusUnauthorized, []byte(expectedSig), func(wallet *keymanager.KeyManager) {
	//		actualSignature, err := wallet.Sign(context.Background(), testRequest(t))
	//		require.True(t, keymanager.IsGenericError(err))
	//		require.Nil(t, actualSignature)
	//	})
	//})

	//t.Run("rejects with failed", func(t *testing.T) {
	//	runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
	//		actualSignature, err := wallet.SignProposal(domain, data)
	//		require.True(t, keymanager.IsGenericError(err))
	//		require.Nil(t, actualSignature)
	//	})
	//})
	//
	//t.Run("rejects with invalid signature", func(t *testing.T) {
	//	runTest(t, http.StatusOK, []byte("invalid"), func(wallet *keymanager.KeyManager) {
	//		actualSignature, err := wallet.SignProposal(domain, data)
	//		require.True(t, keymanager.IsGenericError(err))
	//		require.Nil(t, actualSignature)
	//	})
	//})
}

func testRequest(t *testing.T) *v2.SignRequest {
	// block
	blockByts := "7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d"
	blk := &eth.BeaconBlock{}
	require.NoError(t, json.Unmarshal(_byteArray(blockByts), blk))

	// build request
	return &v2.SignRequest{
		PublicKey:       _byteArray("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972"),
		SigningRoot:     nil,
		SignatureDomain: _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
		Object:          &v2.SignRequest_Block{Block: blk},
	}
}
