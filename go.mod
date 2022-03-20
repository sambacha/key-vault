module github.com/bloxapp/key-vault

go 1.15

require (
	github.com/bloxapp/eth2-key-manager v1.1.1
	github.com/docker/docker v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/getsentry/sentry-go v0.11.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/vault/api v1.1.1
	github.com/hashicorp/vault/sdk v0.2.1
	github.com/herumi/bls-eth-go-binary v0.0.0-20210917013441-d37c07cfda4e
	github.com/makasim/sentryhook v0.4.0
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/eth2-types v0.0.0-20210303084904-c9735a06829d
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210809151128-385d8c5e3fb7
	github.com/prysmaticlabs/prysm v1.4.2-0.20210901151710-46b22bb64999
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494
	google.golang.org/grpc v1.40.0
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20210707101027-e8523651bf6f

// go.mode throwing error when prysm point dirctly to v2.X so replace to v2.0.1 here
replace github.com/prysmaticlabs/prysm => github.com/prysmaticlabs/prysm v1.4.2-0.20220124113610-e26cde5e091b

replace github.com/json-iterator/go => github.com/prestonvanloon/go v1.1.7-0.20190722034630-4f2e55fcf87b

// See https://github.com/prysmaticlabs/grpc-gateway/issues/2
replace github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/prysmaticlabs/grpc-gateway/v2 v2.3.1-0.20210702154020-550e1cd83ec1

replace github.com/ferranbt/fastssz => github.com/prysmaticlabs/fastssz v0.0.0-20220110145812-fafb696cae88
