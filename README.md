

https://developer.android.com/develop/connectivity/nfc/nfc

https://learn.adafruit.com/adafruit-pn532-rfid-nfc/about-nfc
https://cdn-shop.adafruit.com/datasheets/pn532um.pdf
https://en.wikipedia.org/wiki/ISO/IEC_14443
https://www.pjrc.com/store/teensy32.html
https://github.com/don/NDEF
https://github.com/Seeed-Studio/PN532
https://developer.android.com/develop/connectivity/nfc/hce


```
message=NdefRecord tnf=1 type=55 payload=00636F76656E3A2F2F746865636F76656E2E73706163652F646F6F72 id= payload-ascii=��coven://thecoven.space/door
message=NdefRecord tnf=4 type=616E64726F69642E636F6D3A706B67 payload=73706163652E746865636F76656E2E616E64726F6964 id= payload-ascii=space.thecoven.android
```




```
Found an ISO14443A card
  UID Length: 4 bytes
  UID Value: 94 7A 1F 46 
Seems to be a Mifare Classic card (4 byte UID)
------------------------Sector 0-------------------------
Block 0   94 7A 1F 46 B7 08 04 00 04 1A 97 42 4F A7 0F 90    .z.F.......BO...
Block 1   14 01 03 E1 03 E1 03 E1 03 E1 03 E1 03 E1 03 E1    ................
Block 2   03 E1 03 E1 03 E1 03 E1 03 E1 03 E1 03 E1 03 E1    ................
Block 3   00 00 00 00 00 00 78 77 88 C1 00 00 00 00 00 00    ......xw........
------------------------Sector 1-------------------------
Block 4   03 13 D1 01 0F 55 03 74 68 65 63 6F 76 65 6E 2E    .....U.thecoven.
Block 5   73 70 61 63 65 FE 00 00 00 00 00 00 00 00 00 00    space...........
Block 6   00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
Block 7   00 00 00 00 00 00 7F 07 88 40 00 00 00 00 00 00    ........@......
------------------------Sector 2-------------------------
Block 8   00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
Block 9   00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
Block 10  00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00    ................
Block 11  00 00 00 00 00 00 7F 07 88 40 00 00 00 00 00 00    ........@......
...
```



```C
// nfc.InJumpForDEP(); - Active, baud? PassiveInitiatorData, NFCID3i, Gi (ATR_REQ)
nfc.inListPassiveTarget(); // returns true when a card is found (or alternatively readPassiveTargetID)
nfc.inDataExchange(send_buf, send_len, recv_buf, *recv_len); // transfer data via ISO-14443-3A
// when tg is ISO/IEC14443-4 compliant, 
    // -> D4 40 01  [00 B0 82 00 10]
    // <- D5 41 00  [00 01 02 03 04 … 0F 90 00]


// InSelect - trigger initialization - is this necessary??
```

pcd=reader AKA
picc=card  AKA


In=initiator
Tg=target


passive=initiator has radio
active=both sides have radios


ISO-14443-4 reader/card data exchange protocol
NFC-DEP active peer-to-peer  LLCP/SNEP



The Initiator selects the communication mode (either Active or Passive) and bit rate


Android apps use DEP (active extension of 14443-4), on top of ISO-14443-3A

reader: -> "SELECT AID" APDU


If you don't want to register an AID, you are free to use AIDs in the proprietary range: bits 8-5 of the first byte must each be set to '1'. For example, "0xF00102030405" is a proprietary AID.

FF + thecoven.space
FF 74 68 65 63 6F 76 65 6E 2E 73 70 61 63 65 FF


CLA As defined in 5.1.1
INS 'A4'
P1 See Table 39
P2 See Table 40
L c field Absent for encoding N c = 0, present for encoding N c > 0
Data field Absent or file identifier or path or DF name (according to P1)
L e field Absent for encoding N e = 0, present for encoding N e > 0
Data field Absent or file control information (according to P2)
SW1-SW2 See Tables 5 and 6 when relevant, e.g., '6283', '6284', '6A80', '6A81', '6A82', '6A86', '6A87'


```C
cla = 0; // The command is the last or only command of a chain,
        // No SM or no indication, Logical channel number from zero to three
ins = 0xA4; // SELECT
param[0] = 0x04; // the command data contains a DF name (the AID)
param[0]
tag = 0x4F; // AID 8.2.1.2
```


Two categories of structures are supported: dedicated file (DF) and elementary file (EF)
An internal EF stores data interpreted by the card, i.e., data used by the card for management and
control purposes.
Any appli-
cation identifier (AID, see 8.2.1.2) may be used as DF name.

The entity to authenticate has to prove the knowledge of the relevant
secret or private key in an authentication procedure (e.g., a GET CHALLENGE command followed by an
EXTERNAL AUTHENTICATE command, a sequence of GENERAL AUTHENTICATE commands).



GENERAL AUTHENTICATE
INS= '86' or '87
p1=algorithm (or zero, no info given)
p2=secret id (or zero, no info given)




The historical bytes indicate operating characteristics of the card.
The first historical byte is the “category indicator byte”. If the category indicator byte is set to '00', '10' or '8X',
then Table 83 summarizes the format of the historical bytes. Any other value indicates a proprietary format

AID
Referenced by a compact header set to 'FY' in the historical bytes (see 8.1.1), or by tag '4F' in the initial dat
