#include <unity.h>

#include "key_verification.h"

KeyVerification verifier;

void setUp(void) {
    // set stuff up here
}

void tearDown(void) {
    // clean stuff up here
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
    for (size_t i = 0; i < CHALLENGE_SIZE; i++) {
        if (challenge[i] < 16) Serial.print("0");
        Serial.print(challenge[i], HEX);
    }
    Serial.println();
}

// void test_rsa_keygen() {
//     char pubkey[256];
//     memset(pubkey, 0, 256);
//     size_t* outLen;
//     int res = verifier.getEncryptionPublicKey(pubkey, 256, outLen);
//     TEST_ASSERT_EQUAL_INT(0, res);
//     TEST_ASSERT_EQUAL_INT(140, outLen);
//     Serial.println(String(pubkey));
// }

int main( int argc, char **argv) {
    UNITY_BEGIN();

    RUN_TEST(test_begin);
    // RUN_TEST(test_generate_challenge);
    // RUN_TEST(test_rsa_keygen);


    UNITY_END();
}
