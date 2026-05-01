#include <unity.h>

#include "key_verification.h"

KeyVerification verifier;

void setUp(void) {
    // set stuff up here
}

void tearDown(void) {
    // clean stuff up here
}

void debugPrintHex(uint8_t* ptr, size_t size) {
    for (size_t i = 0; i < size; i++) {
        if (ptr[i] < 16) Serial.print("0");
        Serial.print(ptr[i], HEX);
    }
    Serial.println();
}


void test_begin() {
    int res = verifier.begin();
    TEST_ASSERT_EQUAL_INT(0, res);
}

void test_generate_challenge() {
    uint8_t challenge[CHALLENGE_SIZE];
    memset(challenge, 0, CHALLENGE_SIZE);
    int res = verifier.generateChallenge(challenge);
    TEST_ASSERT_EQUAL_INT(0, res);
    Serial.print("Challenge=");
    debugPrintHex(challenge, CHALLENGE_SIZE);
}

void test_generate_exchange_key() {
    uint8_t* key;
    uint16_t outLen = 0;
    int res = verifier.getEncryptionPublicKey(&key, &outLen);
    TEST_ASSERT_EQUAL_INT(0, res);
    TEST_ASSERT_EQUAL_INT(65, outLen);
    Serial.print("ExchangeKey=");
    debugPrintHex(key, (size_t) outLen);
}

void test_decode_key() {
    uint8_t user_secret[96] = {
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
        0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
        0x46, 0xca, 0xa3, 0xf7, 0x0e, 0xed, 0x84, 0x14, 0x40, 0xc8, 0xb3, 0x08, 0xd5, 0xb8, 0xaa, 0x76,
        0xad, 0xf9, 0x88, 0x7d, 0x29, 0xcf, 0xf1, 0x1e, 0x09, 0x60, 0xbb, 0xfe, 0x15, 0xc0, 0xf5, 0x0b,
        0xef, 0x02, 0xa8, 0x25, 0x01, 0x1a, 0xab, 0x88, 0x33, 0x75, 0xd0, 0x15, 0xf4, 0x1a, 0xed, 0x0e,
        0x85, 0x75, 0x9d, 0xaa, 0x97, 0xb2, 0x02, 0x04, 0x2f, 0x77, 0x1e, 0x17, 0xc8, 0x76, 0xac, 0x0f,        
    };
    uint8_t access_key[256];
    int res = verifier.test_generateFakeAccessKey(user_secret, 96, access_key);
    TEST_ASSERT_EQUAL_INT(0, res);
    debugPrintHex(access_key, 256);

    res = verifier.verifyAccessKey(access_key, 256);
    TEST_ASSERT_EQUAL_INT(0, res);
}

int main( int argc, char **argv) {
    UNITY_BEGIN();

    RUN_TEST(test_begin);
    RUN_TEST(test_generate_challenge);
    RUN_TEST(test_generate_exchange_key);
    RUN_TEST(test_decode_key);

    UNITY_END();
}
