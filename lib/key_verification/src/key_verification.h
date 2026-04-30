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

class KeyVerification {
    private:
        WC_RNG kv_rng;
        // ed25519_key server_pub_key;
        ed25519_key door_sign_key;
        Hpke app_exchange;

    public:
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

        wc_HpkeInit(&app_exchange, DHKEM_P256_HKDF_SHA256, HKDF_SHA256, HPKE_AES_128_GCM);

        return 0;
    }

    // generate a challenge, which is a 64-bit secure-random nonce followed
    // by a signature signed by door_sign_priv_key using ed25519. Signature is 64 bytes,
    // meaning challenge must be 8+64=72 bytes.
    int generateChallenge(uint8_t* challenge) {
        int ret = wc_RNG_GenerateBlock(&kv_rng, challenge, NONCE_SIZE);
        if (ret != 0) return ret;

        unsigned int outLen = ED25519_SIG_SIZE;
        ret = wc_ed25519_sign_msg(challenge, NONCE_SIZE,
            challenge+NONCE_SIZE, &outLen, &door_sign_key);
        if (ret != 0) return ret;
        return 0;
    }

    // int getEncryptionPublicKey(uint8_t* output, uint32_t inLen, size_t* outLen) {
    //     int res = wc_RsaKeyToPublicDer_ex(&door_enc_key, output, inLen, false);
    //     if (res < 0) return res;
    //     *outLen = res;
    //     return 0;
    // }

    // int verifyAccessKey(uint8_t* accessKey) {

    // }
};

#endif // KEY_VERIFICATION_H_
