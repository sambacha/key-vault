# Ability to list existing wallet accounts ("list")
path "ethereum/+/accounts" {
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

# Ability to update storage ("create")
path "ethereum/+/storage" {
  capabilities = ["create"]
}

# Ability to read slashing storage ("read")
path "ethereum/+/storage/slashing" {
  capabilities = ["read"]
}
