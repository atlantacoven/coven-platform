#include <types.h>




struct APDU {
    uint8_t cla; // class
    uint8_t ins; // instruction
    uint8_t param[2]; // parameter bytes 1 and 2
    uint8_t* data;
    size_t len;
    bool has_response;
};

size_t write_apdu(APDU* msg, uint8_t* out) {
    size_t i = 0;
    out[i++] = msg->cla;
    out[i++] = msg->ins;
    out[i++] = msg->param[0];
    out[i++] = msg->param[1];
    if (msg->data != NULL) {
        // TODO: support extended:
        // Nc denotes the number of bytes in the command data field. The Lc field encodes Nc.
        // If the Lc field is absent, then Nc is zero.
        // A short Lc field consists of one byte not set to '00'.
        // From '01' to 'FF', the byte encodes Nc from one to 255.
        // An extended Lc field consists of three bytes: one byte set to '00' followed by two bytes not set to '0000'.
        // From '0001' to 'FFFF', the two bytes encode Nc from one to 65_535

        // exclude Lc if no data
        out[i++] = len;
        for (size_t j = 0; j < min(255, len); j++) {
            out[i++] = data[j];
        }
    }
    // Ne denotes the maximum number of bytes expected in the response data field. The Le field encodes Ne.
    // If the Le field is absent, then Ne is zero.
    // A short Le field consists of one byte with any value.
    // From '01' to 'FF', the byte encodes Ne from one to 255.
    // If the byte is set to '00', then N e is 256.
    // An extended Le field consists of either three bytes (one byte set to '00' followed by two bytes with any
    //   value) if the Lc field is absent, or two bytes (with any value) if an extended Lc field is present.
    // From '0001' to 'FFFF', the two bytes encode Ne from one to 65_535.
    // If the two bytes are set to '0000', then Ne is 65_536.
    if (msg->has_response) {
        out[i++] = 0xFF;
    }
    return i;
}

struct SimpleTLV {
    uint8_t tag;
    uint8_t* data;
    uint8_t len; // no more than 254
};


struct Response {
    uint8_t* data;
    size_t len;
    uint16_t status; // 6XXX,9XXX
};
