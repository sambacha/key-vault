package keymanager

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/prysmaticlabs/prysm/shared/bls"

	"github.com/bloxapp/key-vault/backend"
)

// SignGeneric implements ProtectingKeyManager interface.
func (km *KeyManager) SignGeneric(root []byte, domain [32]byte) (bls.Signature, error) {
	// Prepare request body.
	req := SignAggregationRequest{
		PubKey:     km.originPubKey,
		Domain:     hex.EncodeToString(domain[:]),
		DataToSign: hex.EncodeToString(root[:]),
	}

	// Json encode the request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, NewGenericError(err, "failed to marshal request body")
	}

	// Send request.
	var resp SignResponse
	if err := km.sendRequest(http.MethodPost, backend.SignAggregationPattern, reqBody, &resp); err != nil {
		km.log.WithError(err).Error("failed to send sign aggregation request")
		return nil, NewGenericError(err, "failed to send SignGeneric request to remote vault wallet")
	}

	// Signature is base64 encoded, so we have to decode that.
	decodedSignature, err := hex.DecodeString(resp.Data.Signature)
	if err != nil {
		return nil, NewGenericError(err, "failed to base64 decode")
	}

	// Get signature from bytes
	sig, err := bls.SignatureFromBytes(decodedSignature)
	if err != nil {
		return nil, NewGenericError(err, "failed to get BLS signature from bytes")
	}

	return sig, nil
}
