#!/bin/bash

set -eo pipefail

export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)

# Register plugin
vault plugin register \
    -sha256=${SHASUM256} \
    -args="--log-format=${LOG_FORMAT}" \
    -args="--log-dsn=${LOG_DSN}" \
    -args="--log-levels=${LOG_LEVELS}" \
    secret ethsign

# Enable secrets
if vault secrets tune ethereum/pyrmont/ ; then
  echo "the secret at path ethereum/pyrmont already exists"
else
  echo "Enabling pyrmont plugin..."
  vault secrets enable \
    -path=ethereum/pyrmont \
    -description="Eth Signing Wallet - Pyrmont Test network" \
    -plugin-name=ethsign plugin
  echo "Enabled plugin"
fi

if vault secrets tune ethereum/mainnet/ ; then
  echo "the secret at path ethereum/mainnet already exists"
else
  echo "Enabling mainnet plugin..."
  vault secrets enable \
    -path=ethereum/mainnet \
    -description="Eth Signing Wallet - Mainnet network" \
    -plugin-name=ethsign plugin
  echo "Enabled plugin"
fi

# Configuring networks
echo "Configuring Pyrmont Test network..."
vault write ethereum/pyrmont/config \
    network="pyrmont"
echo "Configured Pyrmont Test network"

echo "Configuring MainNet network..."
vault write ethereum/mainnet/config \
    network="mainnet"
echo "Configured MainNet network"

TOKEN=$(cat /data/keys/vault.root.token)

# Reload plugin
curl --insecure --header "X-Vault-Token: $TOKEN" \
     --request PUT \
     --data '{"plugin": "ethsign"}' \
     ${VAULT_SERVER_SCHEMA:-http}://127.0.0.1:8200/v1/sys/plugins/reload/backend

# Make sure everything works well
curl --insecure \
     --header "X-Vault-Token: $TOKEN" \
     --request GET \
     --fail \
     ${VAULT_SERVER_SCHEMA:-http}://127.0.0.1:8200/v1/ethereum/pyrmont/config

curl --insecure \
     --header "X-Vault-Token: $TOKEN" \
     --request GET \
     --fail \
     ${VAULT_SERVER_SCHEMA:-http}://127.0.0.1:8200/v1/ethereum/mainnet/config
