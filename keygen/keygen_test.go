package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func TestGeneratedSigningKey(t *testing.T) {
	pubkey := must(hex.DecodeString("18dda9e0ea95f2b07d62f8146f8b3e02f7fe96c0fc54552407478e6c908a4951"))
	challenge := must(hex.DecodeString("8E345E9BDB9E4843B83CEA881EA365399B563FB307F1E8A43014A65FFD892C0ACCE922C9E0D306CD5F8E34F4CD64E1DACCEAB605D4D05366A81A0292AE890943CD7340AB2F69B103"))

	key := ed25519.PublicKey(pubkey)

	nonce := challenge[0:8]
	sig := challenge[8:]

	assert.Truef(t, ed25519.Verify(key, nonce, sig), "signature valid")
}
