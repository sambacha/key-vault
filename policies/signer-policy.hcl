# Ability to list prater wallet accounts ("list")
path "ethereum/prater/accounts" {
  capabilities = ["list"]
}

# Ability to list mainnet wallet accounts ("list")
path "ethereum/mainnet/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/+/accounts/sign" {
  capabilities = ["create"]
}

# Ability to get version ("read")
path "ethereum/+/version" {
  capabilities = ["read"]
}
