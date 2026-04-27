#include <Arduino.h>
#include <SPI.h>
#include <PN532_SPI.h>
#include <PN532_SPI.cpp>
#include "PN532.h"

#define LED 13
#define NFC_CS_PIN 15

PN532_SPI pn532spi(SPI, NFC_CS_PIN);
PN532 nfc(pn532spi);

uint8_t messagebuf[256];
size_t messageSize;

void setup() {
  Serial.begin(115200);

  nfc.begin();
  size_t version = nfc.getFirmwareVersion();
  if (version == 0) {
    Serial.println("PN532 not found");
    while (1) ;;
  }
  Serial.print("Firmware version: "); Serial.println(version);
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
  0xFA, // proprietary cla value (bit8=1)
  0x90, 0x00, // status ok
};

uint8_t FAKE_KEY[] = { 0xDE, 0xAD, 0xBE, 0xEF };

uint8_t checkKey(uint8_t* key, uint8_t len) {
  if (len != 4 + 2) return 0x10; // invalid length
  uint16_t status;
  status = (key[4] << 8) + key[5];
  if (status != 0x9000) {
    return 0x60; // invalid response status
  }
  for (size_t i = 0; i < 4; i++) {
    if (key[i] != FAKE_KEY[i]) return 0x20; // invalid key
  }
  return 0; // valid
}

void loop() {
    // wait for a card in range
    if (!nfc.inListPassiveTarget()) {
        delay(100);
        return;
    }

    uint8_t recv_len;
    if (!nfc.inDataExchange(AID_MESSAGE, 22, messagebuf, &recv_len)) {
        Serial.println("send failed");
    } else {
        Serial.print("response received len="); Serial.println(recv_len);
        for (size_t i = 0; i < recv_len; i++) {
          Serial.print(messagebuf[i], HEX);
        }
        Serial.println();

        uint8_t keyResult = checkKey(messagebuf, recv_len);
        Serial.print("key result="); Serial.println(keyResult);
        if (keyResult == 0) {
          AUTH_RESULT_MESSAGE[1] = 0x90;
          AUTH_RESULT_MESSAGE[2] = 0x00;
          nfc.inDataExchange(AUTH_RESULT_MESSAGE, 3, messagebuf, &recv_len);
        } else {
          AUTH_RESULT_MESSAGE[1] = 0x66;
          AUTH_RESULT_MESSAGE[2] = keyResult;
          nfc.inDataExchange(AUTH_RESULT_MESSAGE, 3, messagebuf, &recv_len);
        }
        Serial.println("sent response");
    }
    delay(1000);
}
