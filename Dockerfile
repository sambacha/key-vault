FROM buildpack-deps:bullseye-scm@sha256:ecad7c5bf28a62451725306f097e8c731eeeed1c21f4763c1f47eed26e38a9c7 AS preparer

ENV DEBIAN_FRONTEND=noninteractive
ENV GOLANG_VERSION 1.19

ONBUILD ARG GOPROXY
ONBUILD ARG GONOPROXY
ONBUILD ARG GOPRIVATE
ONBUILD ARG GOSUMDB
ONBUILD ARG GONOSUMDB

#RUN GO111MODULE=on GOFLAGS=-mod=vendor go mod vendor
#RUN GO111MODULE=on GOFLAGS=-mod=vendor go mod tidy

#RUN GO111MODULE=on GOFLAGS=-mod=vendor CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#    go build -o $GOLANG_APPLICATION ./cmd/$GOLANG_APPLICATION/entrypoint.go
    
    
# GCC for CGO
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -qqy --assume-yes \
     gcc \
     libc6-dev \
     curl git zip unzip wget g++ jq \
     make; \
     rm -rf /var/lib/apt/lists/*
     apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false;


RUN curl -sSL https://go.dev/dl/go$GOLANG_VERSION.src.tar.gz | tar -v -C /usr/src -xz
RUN export PATH=$PATH:/usr/local/go/bin

#RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV PATH /usr/src/go/bin:$PATH


RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Golang" \
      org.label-schema.description="Primitive Container Enviornment" \
      org.label-schema.url="https://getfoundry.sh" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/manifoldfinance/containers.git" \
      org.label-schema.vendor="CommodityStream, Inc." \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"
      
RUN apt-get update                                                        && \
  DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
  curl git zip unzip wget g++ python gcc jq                                \
  && rm -rf /var/lib/apt/lists/*


WORKDIR /go/src/github.com/bloxapp/key-vault/
COPY go.mod .
COPY go.sum .
RUN go mod download

#
# STEP 2: Build executable binary
#
FROM preparer AS builder

# Copy files and install app
COPY . .
RUN go get -d -v ./...
RUN CGO_CFLAGS="-O -D__BLST_PORTABLE__" CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags \"-static -lm\" -X main.Version=$(git describe --tags $(git rev-list --tags --max-count=1))" -o ethsign .

#
# STEP 3: Get vault image and copy the plugin
#
FROM vault:1.8.12 AS runner

# Download dependencies
RUN apk -v --update --no-cache add \
  bash ca-certificates curl openssl

WORKDIR /vault/plugins/

COPY --from=builder /go/src/github.com/bloxapp/key-vault/ethsign ./ethsign
COPY ./config/vault-config.json /vault/config/vault-config.json
COPY ./config/vault-config-tls.json /vault/config/vault-config-tls.json
COPY ./config/entrypoint.sh /vault/config/entrypoint.sh
COPY ./config/vault-tls.sh /vault/config/vault-tls.sh
COPY ./config/vault-init.sh /vault/config/vault-init.sh
COPY ./config/vault-unseal.sh /vault/config/vault-unseal.sh
COPY ./config/vault-policies.sh /vault/config/vault-policies.sh
COPY ./config/vault-plugin.sh /vault/config/vault-plugin.sh
COPY ./policies/signer-policy.hcl /vault/policies/signer-policy.hcl

RUN chown vault /vault/config/entrypoint.sh
RUN apk add jq

WORKDIR /

# Expose port 8200
EXPOSE 8200

# Run vault
CMD ["bash", "/vault/config/entrypoint.sh"]
