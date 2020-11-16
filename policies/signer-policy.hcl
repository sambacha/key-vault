# Ability to list existing wallet accounts ("list")
path "ethereum/test/accounts" {
  capabilities = ["list"]
}
path "ethereum/mainnet/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/test/accounts/sign-*" {
  capabilities = ["create"]
}
path "ethereum/mainnet/accounts/sign-*" {
  capabilities = ["create"]
}
