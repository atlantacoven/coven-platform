package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/hpke"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

const UserSecretSize = 32
const NonceSize = 8

var ServerKey ed25519.PrivateKey
var ServerPubKey ed25519.PublicKey
var DoorSignKey ed25519.PrivateKey
var DoorSignPubKey ed25519.PublicKey

var KEM = hpke.DHKEM(ecdh.P256())
var KDF = hpke.HKDFSHA256()
var AEAD = hpke.AES128GCM()

var Info = []byte("thecoven.space")

var DoorEncKey hpke.PrivateKey
var DoorEncPubKey []byte

// server has public/private key pair
// door has public/private key pair
// app knows door public key
// door knows server public key
func init() {
	var err error
	ServerPubKey, ServerKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicf("generate key: %v\n", err)
	}
	ServerKeyDer, err := x509.MarshalPKCS8PrivateKey(ServerKey)
	if err != nil {
		log.Panicf("marshall key: %v\n", err)
	}
	fmt.Printf("ServerPrivateKey=%x\n", ServerKeyDer)
	fmt.Printf("ServerPublicKey=%x\n\n", ServerPubKey)

	DoorSignPubKey, DoorSignKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicf("generate key: %v\n", err)
	}
	DoorKeyDer, err := x509.MarshalPKCS8PrivateKey(DoorSignKey)
	if err != nil {
		log.Panicf("marshall key: %v\n", err)
	}
	fmt.Printf("DoorSignPrivateKey=%x\n", DoorKeyDer)
	fmt.Printf("DoorSignPublicKey=%x\n\n", DoorSignKey)

	DoorEncKey, err = KEM.GenerateKey()
	if err != nil {
		log.Panicf("generate key: %v\n", err)
	}
	DoorEncPubKey = DoorEncKey.PublicKey().Bytes()
	DoorEncKeyBytes, err := DoorEncKey.Bytes()
	if err != nil {
		log.Panicf("generate key: %v\n", err)
	}
	fmt.Printf("DoorEncKey[%v]=%x\n", len(DoorEncKeyBytes), DoorEncKeyBytes)
	fmt.Printf("DoorEncPublicKey[%v]=%x\n\n", len(DoorEncPubKey), DoorEncPubKey)
}

func ServerGenerateUserKey(userid uint64) []byte {
	UserSecret := make([]byte, UserSecretSize+ed25519.SignatureSize)
	binary.BigEndian.PutUint64(UserSecret, userid)
	// fill with constant
	for i := 8; i < UserSecretSize; i++ {
		UserSecret[i] = 0xAA
	}

	fmt.Printf("UserSecretDecrypted[%v]=%x\n", UserSecretSize, UserSecret[:UserSecretSize])

	sig := ed25519.Sign(ServerKey, UserSecret[:UserSecretSize])
	copy(UserSecret[UserSecretSize:], sig)
	fmt.Printf("UserSecret[%v]=%x\n\n", len(UserSecret), UserSecret)
	return UserSecret
}

func DoorGenerateChallenge() []byte {
	Nonce := time.Now().Unix()
	fmt.Printf("Nonce=%x\n", Nonce)
	Challenge := make([]byte, NonceSize+ed25519.SignatureSize)
	binary.BigEndian.PutUint64(Challenge, uint64(Nonce))

	DoorSignature := ed25519.Sign(DoorSignKey, Challenge[0:8])
	copy(Challenge[8:], DoorSignature)
	fmt.Printf("Challenge[%v]=%x\n", len(Challenge), Challenge)
	return Challenge
}

func AppGenerateKey(userid uint64, usersecret []byte, challenge []byte, doorSignPubKey ed25519.PublicKey, doorEncPubKey []byte) []byte {
	// verify challenge
	doorSig := challenge[NonceSize:]
	doorSigValid := ed25519.Verify(doorSignPubKey, challenge[0:NonceSize], doorSig)
	if !doorSigValid {
		log.Panicf("Invalid signature")
	}
	nonce := challenge[0:8]

	// generate access key
	AccessKey := make([]byte, NonceSize+len(usersecret))
	copy(AccessKey[0:], nonce)
	copy(AccessKey[NonceSize:], usersecret)
	fmt.Printf("AccessKey[%v]=%x\n", len(AccessKey), AccessKey)

	// encrypt access key with hpke
	publicKey, err := KEM.NewPublicKey(doorEncPubKey)
	if err != nil {
		log.Panicf("get pub key: %v", err)
	}
	FullKey, err := hpke.Seal(publicKey, KDF, AEAD, Info, AccessKey)
	// FullKey is the concatenation of the encapsulated key	and ciphertext
	// For this algorithm, the encapsulated key is 65 bytes
	EncapsulatedKey := FullKey[0:65]
	CipherText := FullKey[65:]

	if err != nil {
		log.Panicf("encrypt key: %v\n", err)
	}
	fmt.Printf("AccessKeyEncapsulatedKey[%v]=%x\n", len(EncapsulatedKey), EncapsulatedKey)
	fmt.Printf("AccessKeyCipherText[%v]=%x\n\n", len(CipherText), CipherText)
	return FullKey
}

func DoorValidateKey(key []byte, serverpubkey ed25519.PublicKey) error {
	DecryptedKey, err := hpke.Open(DoorEncKey, KDF, AEAD, Info, key)
	if err != nil {
		return fmt.Errorf("decrypt key: %w", err)
	}
	fmt.Printf("DecryptedKey=%x\n", DecryptedKey)

	Nonce := DecryptedKey[0:NonceSize]
	fmt.Printf("Nonce=%x\n", Nonce)

	UserSecret := DecryptedKey[NonceSize : NonceSize+UserSecretSize]
	ServerSignature := DecryptedKey[NonceSize+UserSecretSize:]
	fmt.Printf("UserSecret[%v]=%x\n", len(UserSecret), UserSecret)
	fmt.Printf("ServerSignature[%v]=%x\n", len(ServerSignature), ServerSignature)

	valid := ed25519.Verify(serverpubkey, UserSecret, ServerSignature)
	if !valid {
		return fmt.Errorf("invalid UserSecret")
	}

	t := time.Unix(int64(binary.BigEndian.Uint64(Nonce)), 0)
	fmt.Printf("Time=%v\n", t)
	delta := time.Now().Unix() - t.Unix()
	if delta > 15 || delta < -15 {
		return fmt.Errorf("nonce expired")
	}

	srcuid := binary.BigEndian.Uint64(UserSecret[0:8])
	fmt.Printf("UserID=%v\n", srcuid)

	return nil // valid
}

func main() {
	UserId := uint64(0x1234)

	// App authenticates with Server, acquiring UserSecret. It saves this securely to
	// the device
	UserSecret := ServerGenerateUserKey(UserId)

	// When at the door, app generates a temporary Key using the UserSecret, encrypted with
	// the door's public key
	Challenge := DoorGenerateChallenge()
	Key := AppGenerateKey(UserId, UserSecret, Challenge, DoorSignPubKey, DoorEncPubKey)

	// The door validates the key, and if valid unlocks door
	err := DoorValidateKey(Key, ServerPubKey)
	if err != nil {
		log.Panicf("failed: %v", err)
	} else {
		fmt.Printf("okay\n")
	}
}
