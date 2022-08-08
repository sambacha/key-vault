package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bloxapp/key-vault/keymanager/models"

	"github.com/bloxapp/key-vault/utils/encoder/encoderv2"

	bytesutil2 "github.com/prysmaticlabs/prysm/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/runtime/version"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

type signRequestModifier func(block *models.SignRequest)

func signRequestWithFeeRecipient(feeRecipient []byte) signRequestModifier {
	return func(req *models.SignRequest) {
		switch obj := req.Object.(type) {
		case *models.SignRequestBlockV3:
			obj.BlockV3.Body.ExecutionPayload.FeeRecipient = feeRecipient
		}
	}
}

func basicProposalData(blockVersion int, mods ...signRequestModifier) map[string]interface{} {
	return basicProposalDataWithOps(blockVersion, false, false, false, false, mods...)
}

func basicProposalDataWithOps(blockVersion int, undefinedPubKey bool, differentStateRoot bool, differentParentRoot bool, differentBodyRoot bool, mods ...signRequestModifier) map[string]interface{} {
	req := &models.SignRequest{
		PublicKey:       _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
		SigningRoot:     nil,
		SignatureDomain: _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
		Object:          basicSignRequestBlock(blockVersion, differentStateRoot, differentParentRoot, differentBodyRoot),
	}

	if undefinedPubKey {
		req.PublicKey = _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcd")
	}

	for _, mod := range mods {
		mod(req)
	}

	byts, _ := encoderv2.New().Encode(req)
	return map[string]interface{}{
		"sign_req": hex.EncodeToString(byts),
	}
}

func basicSignRequestBlock(blockVersion int, differentStateRoot, differentParentRoot, differentBodyRoot bool) models.ISignObject {
	switch blockVersion {
	case version.Phase0:
		block := &eth.BeaconBlock{}
		jsonData, _ := hex.DecodeString("7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d")
		if err := json.Unmarshal(jsonData, block); err != nil {
			panic(err)
		}

		if differentStateRoot {
			block.StateRoot = _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
		}
		if differentParentRoot {
			block.ParentRoot = _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
		}
		if differentBodyRoot {
			block.Body.Graffiti = bytesutil2.Bytes32(10)
		}

		return &models.SignRequestBlock{Block: block}

	case version.Bellatrix:
		block := &eth.BeaconBlockBellatrix{}
		sszData := _byteArray("1027000000000000b2d6000000000000a665a45cadc0607d6443f800b52c1093b90f88880293dbfefc47579490765fd3dc77754974a255d1e94021c323b04a498e3c5898bfd6b57e76c969a195c65d5c5400000096b3e1eb732fc98e131d6e6c0fe7bfc5e13983005007d48ca3d205ef78d847094d03cdfe326efdb0aaddd4802e23e5e605804cc1484cca30dc42afb9d031f92dcb08e393612b42c3af172ded8c1ebffca15ec4595ea903194d31fd320e19561ed70a234731285c6804c2a4f56711ddb8c82c99740f207854891028af34e27e5e00000000000000000000000000000000000000000000000000000000000000000000000000000000707279736d2d67657468000000000000000000000000000000000000000000008001000080010000800100006b0400006b040000fffeffffeef7ffffffffff7ff7ffdfefdefffff7ffffff4ffbffffe7f7ffefffffef75ffeffbffffffeef7fffbfff6fffefffffeffffeff77f7ffffffeffffdba4875770cd59e728c8875969bb65912e371c3a86a1a04876887a8f084333333964789b8d85403269efd624ba2b81d8180dcfc3dd7537c0dc2e7b2854c8f4d76765a7492c3ede2e128a4e4bef528e28a1e6247e7b4be5881f5747bb1976db021f6b0400000c00000001010000f6010000e40000000f270000000000000000000000000000a665a45cadc0607d6443f800b52c1093b90f88880293dbfefc47579490765fd33701000000000000619a41394302dcc6e512eb09d16202bd5b82ea2f158856be4c5a18ad8409594b3801000000000000594a69a5803accb5c54a68943e4477482c2b29357db63685a9b127a495697cfe8e9bbe3d00e63c05a270efd189a6eb8f20afed7ccacb045e5de940992fe0ab81a9a40cbe429929f76fb059df4a3c27e708e22a7849a7a1b725b68ee1546ec604c00403d9197aca568ee6c1d5d2694fb282656f845e8111a6ef6551006181a73bfffffffffffbffffffdfffffffffffff1fe40000000f270000000000001100000000000000a665a45cadc0607d6443f800b52c1093b90f88880293dbfefc47579490765fd33701000000000000619a41394302dcc6e512eb09d16202bd5b82ea2f158856be4c5a18ad8409594b3801000000000000594a69a5803accb5c54a68943e4477482c2b29357db63685a9b127a495697cfeb3dffe029baae6bea92e6eac8a4ecc4cf6dc07977d2f6ce22b05989ec68cf0922bbd30a50e3c156c54ca53d2d8eabb690d1d386fad0e54cda4d78f473a90736bde3d02481bc613779914cdb952b198f41ce4d4c82a1ee7a427d7a4689838e02af3feffedfffffff7ffffffffffdfbffd1fe40000000f270000000000000400000000000000a665a45cadc0607d6443f800b52c1093b90f88880293dbfefc47579490765fd33701000000000000619a41394302dcc6e512eb09d16202bd5b82ea2f158856be4c5a18ad8409594b3801000000000000594a69a5803accb5c54a68943e4477482c2b29357db63685a9b127a495697cfe8f6206001b2b47be99767b3b1fdb2a813959b0b74b36dcce712b76416477b943db291c1b979aefb3f9aeb177b814cff217ac690bd232b77f078b9534c0ceca711de27f3d48359d72f62ac16f1bce6ee2130ad741d8e48270de5cc9c7381cb36cbfffffffffffffffdbffffffffffffff1f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fc01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fc010000")
		if err := block.UnmarshalSSZ(sszData); err != nil {
			panic(err)
		}

		if differentStateRoot {
			block.StateRoot = _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
		}
		if differentParentRoot {
			block.ParentRoot = _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")
		}
		if differentBodyRoot {
			block.Body.Graffiti = bytesutil2.Bytes32(10)
		}

		block.Body.ExecutionPayload.FeeRecipient = _byteArray("6a3f3eE924A940ce0d795C5A41A817607e520520")

		return &models.SignRequestBlockV3{BlockV3: block}
	}

	panic("block version not supported")
}

var testableBlockVersions = []int{
	version.Phase0,
	version.Bellatrix,
}

// withEachBlockVersion runs the given subtest for each block version.
func withEachBlockVersion(t *testing.T, testName string, test func(t *testing.T, blockVersion int)) {
	for _, blockVersion := range testableBlockVersions {
		t.Run(fmt.Sprintf("%s %s", testName, version.String(blockVersion)), func(t *testing.T) {
			test(t, blockVersion)
		})
	}
}

func TestProposalSlashing(t *testing.T) {
	b, _ := getBackend(t)

	withEachBlockVersion(t, "Successfully Sign proposal", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicProposalData(blockVersion)
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	withEachBlockVersion(t, "Sign proposal with undefined account", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req, func(c *Config) {
			c.FeeRecipients["0x95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcd"] = "0x6a3f3ee924a940ce0d795c5a41a817607e520520"
		})

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicProposalDataWithOps(blockVersion, true, false, false, false)
		res, err := b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: account not found")
		require.Nil(t, res)
	})

	withEachBlockVersion(t, "Sign proposal (exactly same), should error under minimal proposal protection", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first proposal
		req.Data = basicProposalData(blockVersion)
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// second proposal
		req.Data = basicProposalData(blockVersion)
		res, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign: slashable proposal (HighestProposalVote), not signing")
		require.Nil(t, res)
	})

	withEachBlockVersion(t, "Sign double proposal(different state root), should error", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first proposal
		req.Data = basicProposalData(blockVersion)
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// second proposal
		req.Data = basicProposalDataWithOps(blockVersion, false, true, false, false)
		_, err = b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: slashable proposal (HighestProposalVote), not signing")
	})

	withEachBlockVersion(t, "Sign double proposal(different parent root), should error", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first proposal
		req.Data = basicProposalData(blockVersion)
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// second proposal
		req.Data = basicProposalDataWithOps(blockVersion, false, false, true, false)
		_, err = b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: slashable proposal (HighestProposalVote), not signing")
	})

	withEachBlockVersion(t, "Sign double proposal(different body root), should error", func(t *testing.T, blockVersion int) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first proposal
		req.Data = basicProposalData(blockVersion)
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// second proposal
		req.Data = basicProposalDataWithOps(blockVersion, false, false, false, true)
		_, err = b.HandleRequest(context.Background(), req)
		require.NotNil(t, err)
		require.EqualError(t, err, "failed to sign: slashable proposal (HighestProposalVote), not signing")
	})
}
