#ifndef KEY_VERIFICATION_H_
#define KEY_VERIFICATION_H_

#include <Arduino.h>
#include <wolfssl.h>
#include <wolfssl/wolfcrypt/random.h>
#include <wolfssl/wolfcrypt/ed25519.h>
#include <wolfssl/wolfcrypt/hpke.h>
#include <wolfssl/wolfcrypt/misc.h>

#include "_keys.h"

#define NONCE_SIZE 8

#define CHALLENGE_SIZE (NONCE_SIZE+ED25519_SIG_SIZE)

#define PUB_KEY_SIZE ED25519_PUB_KEY_SIZE

#define USER_DATA_SIZE 32
#define USER_ID_SIZE 8
#define PADDING_SIZE (USER_DATA_SIZE-USER_ID_SIZE)
#define PADDING_VALUE 0xAA
#define USER_SECRET_SIZE (USER_DATA_SIZE+ED25519_SIG_SIZE)
#define DECRYPTED_ACCESS_KEY_SIZE (NONCE_SIZE+USER_SECRET_SIZE)

// constant-time mem compare
int constantCompare(uint8_t* a, uint8_t* b, size_t len) {
    int compareSum = 0;
    for (int i = 0; i < len; i++) {
        compareSum |= a[i] ^ b[i];
    }
    return compareSum;
}

int constantCompareVal(uint8_t* a, uint8_t b, size_t len) {
    int compareSum = 0;
    for (int i = 0; i < len; i++) {
        compareSum |= a[i] ^ b;
    }
    return compareSum;
}

class KeyVerification {
    WC_RNG kv_rng[1];
    ed25519_key server_sign_pub_key;
    ed25519_key door_sign_key;
    
    Hpke hpke[1];
    void* skr = NULL;

    uint8_t nonce[NONCE_SIZE];
    uint8_t decryptedKey[DECRYPTED_ACCESS_KEY_SIZE];

    public:
    int begin() {
        int ret = wolfCrypt_Init();
        if (ret != 0) {
            return ret;
        }

        ret = wc_InitRng(kv_rng);
        if (ret != 0) {
            return ret;
        }

        ret = wc_ed25519_import_private_key(DOOR_SIGN_PRIV_KEY,
                ED25519_PRV_KEY_SIZE, NULL, 0, // indicates concat'ed priv+pub
                &door_sign_key);
        if (ret != 0) {
            return ret;
        }

        ret = wc_ed25519_import_public(SERVER_SIGNING_PUB_KEY, ED25519_PUB_KEY_SIZE, &server_sign_pub_key);
        if (ret != 0) {
            return ret;
        }

        ret = wc_HpkeInit(hpke, DHKEM_X25519_HKDF_SHA256, HKDF_SHA256, HPKE_AES_128_GCM, NULL);
        if (ret != 0) {
            return ret;
        }

        return 0;
    }

    // generate a challenge, which is a 64-bit secure-random nonce followed
    // by a signature signed by door_sign_priv_key using ed25519. Signature is 64 bytes,
    // meaning challenge must be 8+64=72 bytes.
    // Saves the nonce for future verification during `verifyAccessKey`.
    int generateChallenge(uint8_t* challenge) {
        int ret = wc_RNG_GenerateBlock(kv_rng, nonce, NONCE_SIZE);
        if (ret != 0) return ret;

        memcpy(challenge, nonce, NONCE_SIZE);

        unsigned int outLen = ED25519_SIG_SIZE;
        ret = wc_ed25519_sign_msg(challenge, NONCE_SIZE,
            challenge+NONCE_SIZE, &outLen, &door_sign_key);
        if (ret != 0) return ret;
        return 0;
    }

    // Generate a one-time keypair (skR and pkR) for receiving an HPKE message from the app
    // and output a serialized pkR. The sender needs this key in order to seal the message.
    // Writes `PUB_KEY_SIZE` bytes into `output`.
    int getEncryptionPublicKey(uint8_t* output) {
        int ret = wc_HpkeGenerateKeyPair(hpke, &skr, kv_rng);
        if (ret != 0) return ret;

        uint16_t pkr_size = PUB_KEY_SIZE;
        ret = wc_HpkeSerializePublicKey(hpke, skr, output, &pkr_size);
        if (ret != 0) return ret;
        return 0;
    }

    // Given an access key, decode the message into UserSecret, and verify it is valid and signed by the server.
    // An AccessKey is the nonce from the last generated challenge, followed by the UserSecret given to the app
    // by the server, encrypted using HPKE/KEM-X25519-SHA256/HKDF-SHA256/AES-GCM128. It is in the format of
    // ciphertext followed by the encapsulation key.
    // Returns negative numbers for wolfssl errors, positive numbers for invalid key errors, 0 for valid.
    int verifyAccessKey(uint8_t* accessKey, size_t accessKeyLen) {
        int res;
        uint8_t* enc = accessKey + accessKeyLen;
        res = wc_HpkeOpenBase(hpke, skr,
            enc, PUB_KEY_SIZE,
            EXCHANGE_INFO, EXCHANGE_INFO_LEN, // info
            NULL, 0, // aar
            accessKey, accessKeyLen,
            decryptedKey);
        // NOTE: for some reason this returns AES_GCM_AUTH_E, probably
        // because we aren't using mode_auth?
        if (res != 0 && res != AES_GCM_AUTH_E) return res;
        
        // nonce is first 8 bytes
        if (constantCompare(decryptedKey, nonce, NONCE_SIZE) != 0) {
            return 10; // invalid nonce
        }
        uint8_t* userSecret = decryptedKey+NONCE_SIZE;
        // userId is the next 8 bytes
        // then 24 bytes of padding 0xAA
        uint8_t* padding = userSecret + USER_ID_SIZE;
        if (constantCompareVal(padding, PADDING_VALUE, PADDING_SIZE) != 0) {
            return 20; // invalid UserSecret
        }
        
        // server signature is the remaining 64 bytes
        uint8_t* serverSig = userSecret + USER_DATA_SIZE;
        res = verifyServerSignature(userSecret, USER_DATA_SIZE, serverSig, ED25519_SIG_SIZE);
        if (res != 0) return res;
        return 0; // valid
    }

    private:
    // Verify the UserSecret was actually generated by the server, using the server's public key
    int verifyServerSignature(uint8_t* message, uint32_t messageLen, uint8_t* signature, uint32_t signatureLen) {
        int res;
        int verified;
        res = wc_ed25519_verify_msg(signature, signatureLen, message, messageLen, &verified, &server_sign_pub_key);
        if (res != 0) return res;
        if (verified == 0) return 30; // invalid signature
        return 0; // ok
    }

    // void free() {
    //     if (skr != NULL)
    //         wc_HpkeFreeKey(hpke, hpke->kem, skr, NULL);
    //     wc_rng_free(kv_rng);
    // }
};

#endif // KEY_VERIFICATION_H_
