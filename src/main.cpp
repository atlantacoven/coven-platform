#include <Arduino.h>

#include <SPI.h>
#include <PN532_SPI.h>
#include <PN532_SPI.cpp>
#include "PN532.h"


#include "key_verification.h"

#if defined ARDUINO_ESP32_THING
#define NFC_CS_PIN 2
#elif defined ARDUINO_TEENSY31
#define NFC_CS_PIN 15
#endif

#define STATUS_OK 0x9000
#define DOOR_UNLOCK_RESULT_CMD 0xFA

PN532_SPI pn532spi(SPI, NFC_CS_PIN);
PN532 nfc(pn532spi);

KeyVerification verifier;

// Buffer for storing messages across NFC
uint8_t messagebuf[256];
size_t messageSize;

void setup() {
  Serial.begin(115200);

  nfc.begin();
  size_t version = nfc.getFirmwareVersion();
  if (version == 0) {
    Serial.println(F("PN532 not found"));
    while (1) ;;
  }
  Serial.print(F("Firmware version: ")); Serial.println(version);
  nfc.SAMConfig();
}

uint8_t AID_MESSAGE[22] = {
  0x00, // CLA
  0xA4, // INS: SELECT command
  0x04, // p1: the command data contains a DF name (the AID)
  0x00, // p2
  16,   // Lc (message len)
  // message data: the AID
  0xFF, 't', 'h', 'e', 'c', 'o', 'v', 'e', 'n', '.', 's', 'p', 'a', 'c', 'e', 0xFF,
  0x00, // accept up to 256 bytes in response
};

uint8_t AUTH_RESULT_MESSAGE[3] = {
  DOOR_UNLOCK_RESULT_CMD, // proprietary cla value (bit8=1)
  0x90, 0x00, // status ok
};

uint8_t AUTH_MESSAGE_HEAD[] = {
  0x00, // CLA
  0x86, // INS: GENERAL AUTHENTICATE
  0x00, 0x00, // params
};

bool statusMatch(uint8_t* statusLoc, uint16_t expected) {
  uint16_t actual = ((uint16_t)(statusLoc[0]) << 8) + ((uint16_t) statusLoc[1]);
  return actual == expected;
}

void debugPrintHex(uint8_t* ptr, size_t size) {
    for (size_t i = 0; i < size; i++) {
        if (ptr[i] < 16) Serial.print("0");
        Serial.print(ptr[i], HEX);
        if ((i+1) % 32 == 0) Serial.println();
    }
    Serial.println();
}

void loop() {
    // wait for a card in range
    if (!nfc.inListPassiveTarget()) {
        delay(100);
        return;
    }

    verifier.begin();

    // send SELECT message with AID
    uint8_t recv_len = 255;
    Serial.println(F("sending SELECT"));
    if (!nfc.inDataExchange(AID_MESSAGE, 22, messagebuf, &recv_len)) {
        Serial.println(F("send failed"));
        return;
    }
    debugPrintHex(messagebuf, recv_len);
    if (recv_len != 2 || !statusMatch(messagebuf, STATUS_OK)) {
      Serial.println(F("not ok"));
      return;
    }

    // Send GENERAL AUTHENTICATE message with Challenge and PubKey
    memcpy(messagebuf, AUTH_MESSAGE_HEAD, 4);
    int res = verifier.generateChallenge(messagebuf + 5);
    if (res != 0) {
      Serial.println(F("Generate challenge fail"));
      return;
    }
    Serial.println("Challenge:"); debugPrintHex(messagebuf + 5, CHALLENGE_SIZE);
    uint16_t messageSize = 0;
    res = verifier.getEncryptionPublicKey(messagebuf + 5 + CHALLENGE_SIZE, &messageSize);
    if (res != 0) {
      Serial.println(F("Generate pubkey fail"));
      return;
    }
    Serial.println("pkR:"); debugPrintHex(messagebuf + 5 + CHALLENGE_SIZE, PUB_KEY_SIZE);
    messageSize += CHALLENGE_SIZE;
    if (messageSize >= 256) {
      Serial.println(F("message overflow"));
    }
    messagebuf[4] = (uint8_t) messageSize;
    messagebuf[5+messageSize] = 0x00; // expect up to 256 bytes

    Serial.println(F("sending GENERAL AUTHENTICATE"));
    recv_len = 255;
    if (!nfc.inDataExchange(messagebuf, 6+messageSize, messagebuf, &recv_len)) {
        Serial.println(F("send failed"));
        return;
    }
    debugPrintHex(messagebuf, recv_len);

    // message=cipher[120]+encapsulatedKey[32]+status[2]
    size_t pksSize = PUB_KEY_SIZE;
    size_t accessKeyLen = recv_len - pksSize - 2;
    if (recv_len < 2 || !statusMatch(messagebuf + pksSize + accessKeyLen, STATUS_OK)) {
      Serial.println(F("not ok"));
      return;
    }
    
    Serial.println("AccessKey:"); debugPrintHex(messagebuf, accessKeyLen);
    Serial.println("pkS:"); debugPrintHex(messagebuf + accessKeyLen, PUB_KEY_SIZE);
    res = verifier.verifyAccessKey(messagebuf, accessKeyLen, messagebuf + accessKeyLen);
    if (res == 0) {
      Serial.println(F("verification success"));
      AUTH_RESULT_MESSAGE[1] = 0x90;
      AUTH_RESULT_MESSAGE[2] = 0x00;
      // TODO: unlock door
    } else {
      Serial.print(F("verification failed: "));
      Serial.println(res);
      AUTH_RESULT_MESSAGE[1] = 0x66;
      AUTH_RESULT_MESSAGE[2] = (uint8_t) res;
    }
    nfc.inDataExchange(AUTH_RESULT_MESSAGE, 3, messagebuf, &recv_len);

    // verifier.free();

    delay(1000);
}
