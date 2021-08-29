package tests

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e"

	"github.com/bloxapp/eth2-key-manager/core"
)

// AggregationReferenceSigning tests sign aggregation endpoint.
// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/develop/beacon-chain/sync/validate_aggregate_proof_test.go#L300
type AggregationReferenceSigning struct {
}

// Name returns the name of the test.
func (test *AggregationReferenceSigning) Name() string {
	return "Test aggregation reference signing"
}

// Run run the test.
func (test *AggregationReferenceSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	storage := setup.UpdateStorage(t, core.PyrmontNetwork, true, core.NDWallet, _byteArray("6327b1e58c41d60dd7c3c8b9634204255707c2d12e2513c345001d8926745eea"))
	// Get wallet
	wallet, err := storage.OpenWallet()
	require.NoError(t, err)
	account, err := wallet.AccountByPublicKey("954eb88ed1207f891dc3c28fa6cfdf8f53bf0ed3d838f3476c0900a61314d22d4f0a300da3cd010444dd5183e35a593c")
	require.NoError(t, err)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey()

	// decode object
	agg := &ethpb.AggregateAttestationAndProof{}
	aggAttByts := _byteArray("08161260a1167cdbebeae876b3fa82d4f4c35fc3dc4706c7ae20cee359919fdbc93a2588c3f7a15c80d12a20c78ac6381a9fe35d06f6b8ae7e95fb87fa2195511bd53ce6f385aa71dda52b38771f954348a57acad9dde225da614c50c02173314417b0961ad1010a0109126a1a20eade62f0457b2fdf48e7d3fc4b60736688286be7c7a3ac4c9a16a5e0600bd9e42222122068656c6c6f2d776f726c640000000000000000000000000000000000000000002a221220eade62f0457b2fdf48e7d3fc4b60736688286be7c7a3ac4c9a16a5e0600bd9e41a60b101ab9cd396472716e5334ecbaf797078452117d73596bc5893480ae48f94eee6d5d7dfd67dad69969771f73b75c10816ce412a385cb85cb556d23649d5587cfc7758d95ee5b0ad33ae1a23ecad7fc08a86eba222497d7ed123a46b893393cd")
	require.NoError(t, json.Unmarshal(aggAttByts, agg))
	domain := _byteArray32("060000008c84cda94176cc2b1268357c57c3160131874a4408e155b0db826d11")
	req, err := test.serializedReq(pubKeyBytes, nil, domain, agg)

	// Send sign attestation request
	sig, err := setup.Sign("sign", req, core.PyrmontNetwork)
	require.NoError(t, err)
	expectedSig := _byteArray("8bf29e58a5b594415ce220c3a9f0d64a4cfa44397f92138f8f31849100149e18e0418ed0cb6068f38909b01e9950d7360a8ba1504bd7451c74add42acd82b148ac0b5f3687c429cc571b96307a8902e9976a24747ad68ad21e372302236aab25")
	require.EqualValues(t, expectedSig, sig)
}

func (test *AggregationReferenceSigning) serializedReq(pk, root, domain []byte, agg *ethpb.AggregateAttestationAndProof) (map[string]interface{}, error) {
	req := &validatorpb.SignRequest{
		PublicKey:       pk,
		SigningRoot:     root,
		SignatureDomain: domain,
		Object:          &validatorpb.SignRequest_AggregateAttestationAndProof{AggregateAttestationAndProof: agg},
	}

	byts, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}, nil
}
