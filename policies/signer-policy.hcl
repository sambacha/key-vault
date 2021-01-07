# Ability to list pyrmont wallet accounts ("list")
path "ethereum/pyrmont/accounts" {
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
