#ifndef KEY_VERIFICATION_H_
#define KEY_VERIFICATION_H_

#include <Arduino.h>
#include <wolfssl.h>
#include <wolfssl/wolfcrypt/random.h>
#include <wolfssl/wolfcrypt/ed25519.h>
#include <wolfssl/wolfcrypt/hpke.h>

#include "door_sign_priv_key.h"

#define NONCE_SIZE 8
#define CHALLENGE_SIZE (NONCE_SIZE+ED25519_SIG_SIZE)

#define EXCHANGE_INFO_LEN 14
byte EXCHANGE_INFO[] = "thecoven.space";

class KeyVerification {
    public:
    WC_RNG kv_rng;
    // ed25519_key server_pub_key;
    ed25519_key door_sign_key;
    Hpke app_exchange;
    void* exchange_key = NULL;
    uint8_t nonce[NONCE_SIZE];
    uint8_t exchange_pub_key[128];
    uint16_t exchange_pub_key_size;

    int begin() {
        int ret = wolfCrypt_Init();
        if (ret != 0) {
            return ret;
        }

        ret = wc_InitRng(&kv_rng);
        if (ret != 0) {
            return ret;
        }

        ret = wc_ed25519_import_private_key(DOOR_SIGN_PRIV_KEY, ED25519_PRV_KEY_SIZE,
                NULL, 0, &door_sign_key);
        if (ret != 0) {
            return ret;
        }

        ret = wc_HpkeInit(&app_exchange, DHKEM_P256_HKDF_SHA256, HKDF_SHA256, HPKE_AES_128_GCM, NULL);
        if (ret != 0) {
            return ret;
        }

        return 0;
    }

    // generate a challenge, which is a 64-bit secure-random nonce followed
    // by a signature signed by door_sign_priv_key using ed25519. Signature is 64 bytes,
    // meaning challenge must be 8+64=72 bytes.
    int generateChallenge(uint8_t* challenge) {
        int ret = wc_RNG_GenerateBlock(&kv_rng, nonce, NONCE_SIZE);
        if (ret != 0) return ret;

        memcpy(challenge, nonce, NONCE_SIZE);

        unsigned int outLen = ED25519_SIG_SIZE;
        ret = wc_ed25519_sign_msg(challenge, NONCE_SIZE,
            challenge+NONCE_SIZE, &outLen, &door_sign_key);
        if (ret != 0) return ret;
        return 0;
    }

    int getEncryptionPublicKey(uint8_t* output, uint16_t* outLen) {
        int ret = wc_HpkeGenerateKeyPair(&app_exchange, &exchange_key, &kv_rng);
        if (ret != 0) return ret;

        exchange_pub_key_size = 128;
        ret = wc_HpkeSerializePublicKey(&app_exchange, exchange_key, exchange_pub_key, &exchange_pub_key_size);
        if (ret != 0) return ret;
        memcpy(output, exchange_pub_key, exchange_pub_key_size);
        *outLen = exchange_pub_key_size;

        return 0;
    }

    // returns negative numbers for wolfssl errors, positive numbers for invalid key errors. 0 for valid.
    int verifyAccessKey(uint8_t* accessKey, size_t accessKeyLen) {
        // uint8_t decryptedKey[128];
        // int res = wc_HpkeOpenBase(&app_exchange, exchange_key,
        //     exchange_pub_key, exchange_pub_key_size,
        //     EXCHANGE_INFO, EXCHANGE_INFO_LEN, NULL, 0,
        //     accessKey, accessKeyLen,
        //     decryptedKey);
        // if (res != 0) return res;

        // // free the temporary key
        // wc_HpkeFreeKey(&app_exchange, DHKEM_P256_HKDF_SHA256, exchange_key, NULL);
        return 0;
    }
};

#endif // KEY_VERIFICATION_H_
