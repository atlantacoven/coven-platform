package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

const UserSecretSize = 32

var ServerKey ed25519.PrivateKey
var ServerPubKey ed25519.PublicKey
var DoorKey *rsa.PrivateKey
var DoorPubKey *rsa.PublicKey

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

	DoorKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("generate key: %v\n", err)
	}
	DoorPivKeyDer, err := x509.MarshalPKCS8PrivateKey(DoorKey)
	if err != nil {
		log.Panicf("marshall key: %v\n", err)
	}
	DoorPubKey = &DoorKey.PublicKey
	fmt.Printf("DoorPrivateKey=%x\n", DoorPivKeyDer)
	DoorPubKeyDer := x509.MarshalPKCS1PublicKey(DoorPubKey)
	fmt.Printf("DoorPublicKey=%x\n\n", DoorPubKeyDer)
}

func ServerGenerateUserKey(userid uint64) []byte {
	UserSecret := make([]byte, UserSecretSize+ed25519.SignatureSize)
	binary.BigEndian.PutUint64(UserSecret, userid)
	rand.Read(UserSecret[8:UserSecretSize]) // Salt

	fmt.Printf("UserSecretDecrypted=%x\n", UserSecret)

	// ServerCipher, err := aes.NewCipher(ServerKey)
	// if err != nil {
	// 	log.Panicf("generate cipher: %v\n", err)
	// }
	// for i := 0; i < UserSecretSize; i += ServerCipher.BlockSize() {
	// 	ServerCipher.Encrypt(UserSecret[i:], UserSecret[i:])
	// }
	sig := ed25519.Sign(ServerKey, UserSecret[:UserSecretSize])
	copy(UserSecret[UserSecretSize:], sig)
	fmt.Printf("UserSecret=%x\n\n", UserSecret)
	return UserSecret
}

func AppGenerateKey(userid uint64, usersecret []byte, dpubkey *rsa.PublicKey) []byte {
	Nonce := time.Now().Unix()
	fmt.Printf("Nonce=%x\n", Nonce)

	dkey := make([]byte, 8+UserSecretSize+ed25519.SignatureSize)
	binary.BigEndian.PutUint64(dkey[0:], uint64(Nonce))
	copy(dkey[8:], usersecret)
	fmt.Printf("KeyDecrypted=%x\n", dkey)

	Key, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, dpubkey, dkey, nil)
	if err != nil {
		log.Panicf("encrypt key: %v\n", err)
	}
	fmt.Printf("Key=%x\n\n", Key)
	return Key
}

func DoorValidateKey(key []byte, serverpubkey ed25519.PublicKey) error {
	DecryptedKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, DoorKey, key, nil)
	if err != nil {
		return fmt.Errorf("decrypt key: %w", err)
	}
	fmt.Printf("DecryptedKey=%x\n", DecryptedKey)

	// ServerCipher, err := aes.NewCipher(serverkey)
	// if err != nil {
	// 	return fmt.Errorf("generate cipher: %w", err)
	// }
	// outusec := make([]byte, UserSecretSize)
	// for i := 0; i < UserSecretSize; i += ServerCipher.BlockSize() {
	// 	ServerCipher.Decrypt(outusec[i:], DecryptedKey[16+i:])
	// }
	valid := ed25519.Verify(serverpubkey, DecryptedKey[8:UserSecretSize+8], DecryptedKey[UserSecretSize+8:])
	// fmt.Printf("DecryptedUserSecret=%x\n\n", outusec)
	if !valid {
		return fmt.Errorf("invalid UserSecret")
	}

	t := time.Unix(int64(binary.BigEndian.Uint64(DecryptedKey[0:8])), 0)
	fmt.Printf("Time=%v\n", t)
	delta := time.Now().Unix() - t.Unix()
	if delta > 15 || delta < -15 {
		return fmt.Errorf("nonce expired")
	}

	srcuid := binary.BigEndian.Uint64(DecryptedKey[8:16])
	// destuid := binary.BigEndian.Uint64(DecryptedKey[0:8])
	fmt.Printf("UserID=%v\n", srcuid)
	// if srcuid != destuid {
	// 	return fmt.Errorf("user ids do not match: %v %v", srcuid, destuid)
	// }

	return nil // valid
}

func main() {
	UserId := uint64(0x1234)

	// App authenticates with Server, acquiring UserSecret. It saves this securely to
	// the device
	UserSecret := ServerGenerateUserKey(UserId)

	// When at the door, app generates a temporary Key using the UserSecret, encrypted with
	// the door's public key
	Key := AppGenerateKey(UserId, UserSecret, &DoorKey.PublicKey)

	// The door validates the key, and if valid unlocks door
	err := DoorValidateKey(Key, ServerPubKey)
	if err != nil {
		log.Panicf("failed: %v", err)
	} else {
		fmt.Printf("okay\n")
	}
}
