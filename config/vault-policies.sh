#!/bin/sh

SEALED=$(vault status -format='json' | jq -r '.sealed')
if [[ $SEALED == 'false' ]]; then
 vault policy write signer vault/policies/signer-policy.hcl 2>&1
 vault token create -format=json -policy="signer" > /tmp/vault.signer 2>&1
  cat /tmp/vault.signer | jq -r ".auth.client_token" > /data/keys/vault.signer.token && chmod 0677 /data/keys/vault.signer.token
  rm -f /tmp/vault.*
fi
