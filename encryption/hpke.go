package main

// Try and implement HPKE from scratch, so I can
// properly grok and debug the code

////// KEM: A key encapsulation mechanism

// Ephemerial Key != sender/recepient key!
// Sender/Recepient/Ephemeral

// GenerateKeyPair() -> (skX, pkX)
// SerializePublicKey(pkX) -> []byte
// DeserializePublicKey([]byte) -> pkX

// Generate an ephemeral symmetric key and encapsulation
// which can be de-encapsulated with skR
// Encap(pkR) -> enc:[]byte
// Decap(enc, skR) -> key:[]byte

// Same as encap but proves the sender was the holder of skS
// AuthEncap(pkR, skS) -> enc:[]byte
// AuthDecap(enc, skR, pkS) -> key:[]byte

///// KDF: key derivation

// deterministically generate a key
// Extract(salt?, ikm) -> prk:[]byte
// Expand(prk, info?, L) -> keying_material:[]byte

////// AEAD: encryption algorithm

// Seal(key, nonce, aad?, pt) -> ct:[]byte
// Open(key, nonce, aad?, ct) -> pt:[]byte

// LabeledExtract / LabeledExpand : adds context string to ikm/info

// Encap generates ephemerial keypair
// enc is serialized public key pkE
// kem_context = concat(enc, pkRm)
