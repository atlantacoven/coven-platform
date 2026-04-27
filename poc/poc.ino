#define NDEF_USE_SERIAL 1
#define NFC_INTERFACE_SPI 1

#include <SPI.h>
#include <PN532_SPI.h>
#include <PN532_SPI.cpp>
#include "PN532.h"

#define LED 13
#define NFC_CS_PIN 15

PN532_SPI pn532spi(SPI, NFC_CS_PIN);

PN532 nfc(pn532spi);

const uint8_t AID[16] = { 0xFF, 't', 'h', 'e', 'c', 'o', 'v', 'e', 'n', '.', 's', 'p', 'a', 'c', 'e', 0xFF };

uint8_t messagebuf[128];
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

void loop() {
    // wait for a card in range
    // TODO: probably need to use DEP instead
    if (!nfc.inListPassiveTarget()) {
        delay(100);
        return;
    }
    messagebuf[0] = 0;
    messagebuf[1] = 0xA4; // SELECT
    messagebuf[2] = 0x04; // the command data contains a DF name (the AID)
    messagebuf[3] = 0;
    messagebuf[4] = 16; // the size of the data
    // Data is just the AID (not TLV encoded)
    memcpy(messagebuf+5, AID, 16);

    messagebuf[5+16] = 0xFF; // expect up to 256 bytes back

    for (size_t i = 0; i < 5+16; i++) {
      Serial.print(messagebuf[i], HEX);
    }
    Serial.println();

    uint8_t recv_len;
    if (!nfc.inDataExchange(messagebuf, 8+16, messagebuf, &recv_len)) {
        Serial.println("send failed");
    } else {
        Serial.print("response received len="); Serial.println(recv_len);
        for (size_t i = 0; i < recv_len; i++) {
          Serial.print(messagebuf[i], HEX);
        }
        Serial.println();
    }
    delay(1000);
}