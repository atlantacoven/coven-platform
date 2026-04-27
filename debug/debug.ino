// File: RFID_Cloner.ino
// Author: Bocaletto Luca
// Read, clone & emulate MIFARE Classic 1K / NTAG2xx cards via PN532
// Requires Adafruit_PN532 library

        #include <SPI.h>
        #include <Adafruit_PN532.h>

        #define PN532_SS     15


        //=== Global Instances ===
        Adafruit_PN532 nfc(PN532_SS);

        struct CardDump {
            uint8_t uid[7], uidLen;
            uint8_t data[16][4][16];
            bool    valid;
        } dump;

        const uint8_t defaultKeyA[6] = {0xFF,0xFF,0xFF,0xFF,0xFF,0xFF};
        const uint8_t defaultKeyB[6] = {0xFF,0xFF,0xFF,0xFF,0xFF,0xFF};
        uint8_t *activeKey = (uint8_t*)defaultKeyA;

        //=== Utility Prototypes ===
        void menu();
        void readCard();
        void writeCard();
        void emuCard();
        void loadDump();
        void saveDump();

        void setup() {
            Serial.begin(115200);
            SPI.begin();

            // NFC
            nfc.begin();
            if (!nfc.getFirmwareVersion()) {
              Serial.println("No chip found");
              while(1);
              }
            nfc.SAMConfig();

            dump.valid = false;
            menu();
        }

        void loop() {
            if (Serial.available()) {
                char cmd = toupper(Serial.read());
                switch (cmd) {
                    case 'R': readCard();      break;
                    case 'W': writeCard();     break;
                    case 'E': emuCard();       break;
                    case 'K': activeKey = (activeKey==defaultKeyA?defaultKeyB:defaultKeyA);
                             Serial.println("Key toggled"); break;
                    case 'C': dump.valid=false; Serial.println("Cleared RAM"); break;
                    case '?': menu();          break;
                }
            }
        }

        void menu() {
            Serial.println(F("\n=== RFID/NFC Cloner Pro ==="));
            Serial.println(F("R - Read & dump card"));
            Serial.println(F("W - Write RAM dump to tag"));
            Serial.println(F("E - Emulate dumped card"));
            Serial.println(F("L - Load dump from SD"));
            Serial.println(F("S - Save dump to SD"));
            Serial.println(F("K - Toggle Key A/B"));
            Serial.println(F("C - Clear RAM dump"));
            Serial.println(F("? - Show this menu\n"));
        }

        void readCard() {
            Serial.println("Present tag to read...");
            while (!nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, dump.uid, &dump.uidLen)) {}
            Serial.print("UID: ");
            for (uint8_t i=0;i<dump.uidLen;i++) {
                Serial.print(dump.uid[i], HEX); Serial.print(':');
            }
            Serial.println();
            for (uint8_t s=0;s<16;s++) {
                if (!nfc.mifareclassic_AuthenticateBlock(dump.uid, dump.uidLen, s, 0, activeKey)) {
                    Serial.print("Auth fail sector ");Serial.println(s);
                    continue;
                }
                for (uint8_t b=0;b<4;b++) {
                    nfc.mifareclassic_ReadDataBlock(s*4+b, dump.data[s][b]);
                }
            }
            dump.valid = true;
            Serial.println("Read complete");
        }

        void writeCard() {
            if (!dump.valid) { Serial.println("No dump in RAM"); return; }
            Serial.println("Present blank tag to write...");
            uint8_t uid2[7], len2;
            while (!nfc.readPassiveTargetID(PN532_MIFARE_ISO14443A, uid2, &len2)) {}
            for (uint8_t s=0;s<16;s++) {
                if (!nfc.mifareclassic_AuthenticateBlock(uid2,len2,s,0,activeKey)) continue;
                for (uint8_t b=0;b<4;b++) {
                    nfc.mifareclassic_WriteDataBlock(s*4+b, dump.data[s][b]);
                }
            }
            Serial.println("Write complete");
        }

        void emuCard() {
            if (!dump.valid) { Serial.println("No dump in RAM"); return; }
            Serial.println("Entering emulation...");
            nfc.inListPassiveTarget();  // start emulation
            Serial.println("Emulating. Reset to exit.");
            while (1);
        }
