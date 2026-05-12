package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/hpke"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratedSigningKey(t *testing.T) {
	pubkey := must(hex.DecodeString("18dda9e0ea95f2b07d62f8146f8b3e02f7fe96c0fc54552407478e6c908a4951"))
	challenge := must(hex.DecodeString("8E345E9BDB9E4843B83CEA881EA365399B563FB307F1E8A43014A65FFD892C0ACCE922C9E0D306CD5F8E34F4CD64E1DACCEAB605D4D05366A81A0292AE890943CD7340AB2F69B103"))

	key := ed25519.PublicKey(pubkey)

	nonce := challenge[0:8]
	sig := challenge[8:]

	assert.Truef(t, ed25519.Verify(key, nonce, sig), "signature valid")
}

func TestKeyExchange(t *testing.T) {
	info := []byte("thecoven.space")

	accesskey := must(hex.DecodeString("0000000069f3bd8e0000000000001234aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaacfefdd20526635f07c6c0f2b1627117d6c9e28a5d0648cc4f1f56adc9afcd3b8a078e276cc5a98cfcc89cd9431dcdee9486ef7d46a305e5b5b9efca7e039050a"))

	pubkeydata := must(hex.DecodeString("0448BEB58FF8CC4C2BC2BCC18D8A3F23EB3D94CDCBD64E0EC890D15FE9F71476CFB11F5651A3FD15A03A9E164E5FDDCAE8B9629EF24CCB5146972039F1A25DAC8B"))

	kem := hpke.DHKEM(ecdh.P256())
	kdf := hpke.HKDFSHA256()
	aead := hpke.AES128GCM()

	pubkey := must(kem.NewPublicKey(pubkeydata))
	key := must(hpke.Seal(pubkey, kdf, aead, info, accesskey))

	fmt.Printf("Key[%v]:\n", len(key))
	for _, b := range key {
		fmt.Printf("0x%x, ", b)
	}
	fmt.Println()
}

func TestAndroidOpen(t *testing.T) {
	kem := hpke.DHKEM(ecdh.P256())
	kdf := hpke.HKDFSHA256()
	aead := hpke.AES128GCM()

	info := []byte("thecoven.space")
	cipherText := must(hex.DecodeString("d3881b94d1c4f6bedc901928ec29a841de44bc52ced98c3a751b67d8b80475bc7938fdda69c11b486c7ce85c32e49b5ea704ff5937f4f1fc6402411445aafae93cd2b90470148f39ea9939eed34bbbed1ed1855cdc60d95e1170b6d3f3009c848a4c9038ada7497d82a211df4efe14ba7787fa5f926ec0c8"))
	encapsulatedKey := must(hex.DecodeString("04a4e6072b9e86c97172886586bb35d8ffbed356a804ef70a28cd4d0ab4db8a532cf14df9ebd304aaaf91121c1467d47377d140e798ca0fb984d7175bb30db0a38"))

	privKey := must(kem.NewPrivateKey([]byte(must(hex.DecodeString("ee8834592960d622be5853f7a3441c2e6dad667bd5f7d31a54969cd25d56bb44")))))

	combinedCipher := make([]byte, len(encapsulatedKey)+len(cipherText))
	copy(combinedCipher[0:len(encapsulatedKey)], encapsulatedKey)
	copy(combinedCipher[len(encapsulatedKey):], cipherText)
	decrypted := must(hpke.Open(privKey, kdf, aead, info, combinedCipher))

	expected := must(hex.DecodeString("0000000069f3ee9d0000000000001234aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa23020bc0a02c5b7aa5dd064d2c40ed6996650c58c635597fd9601687c09d7ba7dfba23f585f853d267909efeada28322c3bb280e96d92636a2937c039d922b0a"))
	assert.EqualValues(t, expected, decrypted)
}

func TestDirectOpen(t *testing.T) {
	kem := hpke.DHKEM(ecdh.X25519())
	kdf := hpke.HKDFSHA256()
	aead := hpke.AES128GCM()

	info := []byte("thecoven.space")

	// skR := must(kem.GenerateKey())
	// skR := must(kem.DeriveKeyPair([]byte(must(hex.DecodeString("b8c18a67c46736811bd81c6aad8520d2871cca22844ab4d6883f8c12694f3d4e")))))
	skR := must(kem.NewPrivateKey([]byte(must(hex.DecodeString("d0131ad0e69e97fdc5ea978790f21320bf1a89f2cb7b39a7e84595697ac4fa79")))))
	pkR := skR.PublicKey()
	log.Printf("skR: %x\n", must(skR.Bytes()))
	log.Printf("pkR: %x\n", pkR.Bytes())
	// 5580eb6a4d21e341f82f89b441ae22d18114a123c9e5c64cac9ae6a63d50605f

	ct := []byte(must(hex.DecodeString("9b79ae08c9a84989c4fc0691fa71171ad28291efcab5d646bbe3d4ae0b9b3037ffc02326799364e76a588ac60847c5f65c6d257be8915c555417b88f8c0354338fa8f32fd5d027252ffca598b70168d5c7678f42b51d6d6fe96c28d4663d5b101bb7dbb1c58f0f379bf797c145fe9df536287078a00771c6")))
	enc := []byte(must(hex.DecodeString("5c511df441dd45e854eacea39e565e5dac8a3b9b9e6033026a445c759ca52334")))

	r := must(hpke.NewRecipient(enc, skR, kdf, aead, info))
	pt := must(r.Open(nil, ct))

	log.Printf("pt: %x\n", pt)
}
