# Door Lock

This project implements an electronic NFC door lock for the makerspace from scratch,
allowing active members access to the space via their phone.

It's a large mono-repo consisting of several interconnected subprojects:
- `firmware`: The embedded code for the door lock controller, using PlatformIO
- `server`: An HTTP API in Go for managing and authenticating members
- `android`: An Android app that is able to unlock the door via NFC
- `keygen`: A small Go CLI which generates security credentials for all the projects

## Handshake algorithm

The security requirements for the door lock mechanism:
- An attacker should not be able to generate a valid key in a reasonable time
- The NFC communication should be considered cleartext. An attacker who eavesdropped on a successful NFC handshake should not be able to reuse a previously valid key
- An attacker who spoofs the lock controller should not be able to steal a valid key from the app
- The lock controller should be able to verify authenticity offline, to avoid a network-based attack

Based on these requirements, the following algorithm is used:
1. The app authenticates via HTTPS with the server using a typical mechanism (password, oauth, etc)
2. The server returns a UserSecret, which is info about the user, signed by the server using ED25519.

```
Nonce = random(8)
Padding = repeat(0xAA, 8)
UserData = UserID[8]+ExpiresAt[8]+Nonce[8]+Padding[8]
Signature = ed25519.Sign(UserData, server_private_key)
UserSecret = UserData[32]+Signature[64]
```
3. The app saves this UserSecret securely on the device. It periodically re-authenticates to avoid expiration.
4. The user presents the phone to the door controller. The phone and door controller communicate via standards described in ISO 14443-4 and ISO 7816-4.
5. The door sends an AID identifying it as implementing this (proprietary) algorithm. The app has associated itself with this AID so that it executes in response to the message.
6. The door sends an authentication Challenge to the app, which is a random nonce, a signature verifying itself, and a public key to use for secure transfer of the UserSecret.

```
Nonce = random(8)
Signature = ed25519.Sign(Nonce, door_private_key)
ReceiverKey = hpke.GenerateKeyPair()
ReceiverPublicKey = hpke.SerializePublicKey(ReceiverKey.Public)
Challenge = Nonce[8]+Signature[64]+ReceiverPublicKey[32]
```

7. The app verifies the signature.

```
Nonce = Challenge[0:8]
Signature = Challenge[8:8+64]
ed25519.Verify(Nonce, Signature)
```

8. The app encrypts the UserSecret and sends it. Secure transfer of the UserSecret is done using HPKE as defined in [RFC 9180](https://datatracker.ietf.org/doc/rfc9180/), using `DHKEM(X25519, HKDF-SHA256)/HKDF-SHA256/AES-128-GCM`.

```
ReceiverPublicKey = Challenge[8+64:8+64+32]
ChallengeResponse = Nonce+UserSecret
CipherText, EncapsulationKey = hpke.Seal(ChallengeResponse, ReceiverPublicKey)
EncryptedAccessKey = CipherText[120] + EncapsulationKey[32]
```

9. The door decrypts the message. Then it verifies that the nonce is correct, that the UserSecret is signed by the server, and that it has not expired.

```
CipherText = EncryptedAccessKey[0:120]
EncapsulationKey = EncryptedAccessKey[120:120+32]
ChallengeResponse = hpke.Open(CipherText, EncapsulationKey, ReceiverKey.Private)
UserSecret = ChallengeResponse[8:]
ed25519.Verify(UserSecret[0:32], UserSecret[32:])
equal(Nonce, ChallengeResponse[0:8])
ExpiresAt = UserSecret[8:16]
isBefore(ExpiresAt)
```

10. If all of those checks pass, it unlocks the door for a brief number of seconds. It sends a message to the app indicating if the door was unlocked.

An example of this algorithm can be found in `keygen/algorithm_test.go`.

## Useful Reference

- https://learn.adafruit.com/adafruit-pn532-rfid-nfc/about-nfc
- https://developer.android.com/develop/connectivity/nfc/hce
- http://www.emutag.com/iso/14443-3.pdf
- http://www.emutag.com/iso/14443-4.pdf
- https://www.freecalypso.org/pub/GSM/ISO7816/ISO_7816-4_2005.pdf
- https://cdn-shop.adafruit.com/datasheets/pn532um.pdf
- https://datatracker.ietf.org/doc/rfc9180/
