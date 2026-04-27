#include <Arduino.h>
#include <SPI.h>
#include <PN532_SPI.h>
#include <PN532_SPI.cpp>
#include "PN532.h"
#include "emulatetag.h"
#include <NfcAdapter.h>

#define LED 13
#define NFC_CS_PIN 15

PN532_SPI pn532spi(SPI, NFC_CS_PIN);

// EmulateTag nfc(pn532spi);
// NfcAdapter nfc = NfcAdapter(pn532spi);
PN532 nfc(pn532spi);

const uint8_t uid[3] = {'c', 'v', 'n'};

const uint8_t AID[16] = { 0xFF, 't', 'h', 'e', 'c', 'o', 'v', 'e', 'n', '.', 's', 'p', 'a', 'c', 'e', 0xFF };

uint8_t messagebuf[128];
size_t messageSize;

void setup() {
  Serial.begin(115200);
//   nfc.begin(/*verbose=*/true);

//   NdefMessage message = NdefMessage();
//   // message.addTextRecord("Hello, Arduino!");
//   // message.addUriRecord("coven://thecoven.space/door");
//   message.addUriRecord("http://thecoven.space");

//   NdefRecord appRecord;
//   appRecord.setTnf(TNF_EXTERNAL_TYPE);
//   String appRecordType = "android.com:pkg";
//   String appRecordPayload = "space.thecoven.android";
//   appRecord.setType(appRecordType.c_str(), appRecordType.length());
//   appRecord.setPayload(appRecordPayload.c_str(), appRecordPayload.length());
//   message.addRecord(appRecord);

//   messageSize = message.getEncodedSize();
//   message.encode(messagebuf);
//   Serial.print("set message. bytes="); Serial.println(messageSize);
//   message.print();
//   nfc.setNdefFile(messagebuf, messageSize);
//   nfc.setUid(uid);
//   nfc.init();

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
    messagebuf[4] = 16+2; // the size of the data
    // TLV encoded data
    messagebuf[5] = 0x4F; // AID tag
    messagebuf[6] = 16; // size
    memcpy(messagebuf+7, AID, 16);

    messagebuf[7+16] = 0xFF; // expect up to 256 bytes back

    for (size_t i = 0; i < 8+16; i++) {
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

// void loop() {
//   uint8_t command[] = {
//       PN532_COMMAND_TGINITASTARGET,
//       0x02, // 0x05, // MODE: PICC only, Passive only

//       0x08, 0x00,       // SENS_RES
//       0x12, 0x34, 0x56, // NFCID1
//       0x40,             // SEL_RES
//   };

//   uint8_t status = nfc.tgInitAsTarget(command, 8, 0);
//   if (!status) {
//     Serial.println("Couldn't tginit");
//     while(1);;
//   }

//   nfc.inListPassiveTarget();

//   nfc.setNdefFile(messagebuf, messageSize);

//   nfc.emulate();
//   if (nfc.writeOccured()) {
//     Serial.println("Received write");
//     uint8_t* response_buf;
//     uint16_t messageSize;
//     nfc.getContent(&response_buf, &messageSize);
//     NdefMessage msg = NdefMessage(response_buf, messageSize);
//     msg.print();
//   }
//   delay(1000);

//   if (nfc.tagPresent()) {
//     Serial.println("Saw tag! write message");
//     NdefMessage message = NdefMessage();
//     // message.addTextRecord("Hello, Arduino!");
//     // message.addUriRecord("coven://thecoven.space/door");
//     message.addUriRecord("http://thecoven.space");
    
//     NdefRecord appRecord;

//     appRecord.setTnf(4);
//     String appRecordType = "android.com:pkg";
//     String appRecordPayload = "space.thecoven.android";
//     appRecord.setType(appRecordType.c_str(), appRecordType.length());
//     appRecord.setPayload(appRecordPayload.c_str(), appRecordPayload.length());
//     message.addRecord(appRecord);

//     bool success = nfc.write(message);
//     if (success) Serial.println("message written"); else Serial.println("write failed");
//     delay(1000);
//   }
// }

// void loop() {
//     Serial.println("\nPlace a formatted Mifare Classic or Ultralight NFC tag on the reader.");
//     if (nfc.tagPresent()) {
//         NdefMessage message = NdefMessage();
//         message.addUriRecord("coven://thecoven.space");

//         NdefRecord appRecord;
//         appRecord.setTnf(TNF_EXTERNAL_TYPE);
//         String appRecordType = "android.com:pkg";
//         String appRecordPayload = "space.thecoven.android";
//         appRecord.setType(appRecordType.c_str(), appRecordType.length());
//         appRecord.setPayload(appRecordPayload.c_str(), appRecordPayload.length());
//         message.addRecord(appRecord);

//         bool success = nfc.write(message);
//         if (success) {
//           Serial.println("Success. Try reading this tag with your phone.");        
//         } else {
//           Serial.println("Write failed.");
//         }
//     }
//     delay(5000);
// }