package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const RSA_KEYSIZE = 1024

func main() {
	DoorSignPub, DoorKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicf("gen key: %v\n", err)
	}

	f1, err := os.Create("include/door_sign_priv_key.h")
	if err != nil {
		log.Panicf("create file: %v", err)
	}
	defer f1.Close()
	f1.WriteString("#include <Arduino.h>\n\n")
	fmt.Fprintf(f1, "/* %v */", hex.EncodeToString(DoorKey))
	f1.WriteString("const uint8_t DOOR_SIGN_PRIV_KEY[64] = {\n")
	for _, b := range DoorKey {
		fmt.Fprintf(f1, "0x%x, ", b)
	}
	f1.WriteString("\n};\n")

	// pubkeyhex := hex.EncodeToString(DoorSignPub)
	pubkey, err := x509.MarshalPKIXPublicKey(DoorSignPub)
	pubkeyhex := hex.EncodeToString(pubkey)
	if err != nil {
		log.Panicf("marshall: %v", err)
	}
	err = os.WriteFile("door_signing.pubkey", []byte(pubkeyhex), 0664)
	if err != nil {
		log.Panicf("create file: %v", err)
	}
}
