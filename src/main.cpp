#include <Arduino.h>
#include <SPI.h>
#include <PN532_SPI.h>
#include <PN532_SPI.cpp>
#include "PN532.h"
#include <NfcAdapter.h>

#define LED 13
#define NFC_CS_PIN 15

PN532_SPI pn532spi(SPI, NFC_CS_PIN);
NfcAdapter nfc = NfcAdapter(pn532spi);

void setup() {
  Serial.begin(115200);

  nfc.begin(true /* verbose */);
}

void loop() {
  if (nfc.tagPresent()) {
    Serial.println("Saw tag! write message");
    NdefMessage message = NdefMessage();
    message.addTextRecord("Hello, Arduino!");
    message.addUriRecord("https://example.com");
    bool success = nfc.write(message);
    if (success) Serial.println("message written"); else Serial.println("write failed");
    delay(1000);
  }
}
