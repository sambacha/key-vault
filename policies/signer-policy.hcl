# Ability to list existing wallet accounts ("list")
path "ethereum/pyrmont/accounts" {
  capabilities = ["list"]
}
path "ethereum/mainnet/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/pyrmont/accounts/sign-*" {
  capabilities = ["create"]
}
path "ethereum/mainnet/accounts/sign-*" {
  capabilities = ["create"]
}
