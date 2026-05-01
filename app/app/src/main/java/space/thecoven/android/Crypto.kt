package space.thecoven.android

import android.util.Log
import org.bouncycastle.crypto.hpke.HPKE


fun generateKey() {
//    val ks = KeyStore.getInstance("AndroidKeyStore").also { it.load(null) }
    val info = "thecoven.space".toByteArray(Charsets.US_ASCII)

    val accessKey = "0000000069f3ee9d0000000000001234aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa23020bc0a02c5b7aa5dd064d2c40ed6996650c58c635597fd9601687c09d7ba7dfba23f585f853d267909efeada28322c3bb280e96d92636a2937c039d922b0a".hexToByteArray()

    val pubKeyBytes = "04afc68c08b76f70dec12ccf641929e9bd96129f523ed69e220b37bd3b63e13f142dec9fa4b50b385156bfe123160c35356d127bf9f3289ad484e5bc386c0596ad".hexToByteArray()

    // NOTE: when minSdk is 37 we can use the native Android libs for this instead of bouncycastle
    val hpke = HPKE(HPKE.mode_base, HPKE.kem_P256_SHA256, HPKE.kdf_HKDF_SHA256, HPKE.aead_AES_GCM128)

    val pubKey = hpke.deserializePublicKey(pubKeyBytes)

    val result = hpke.seal(pubKey, info, byteArrayOf(), accessKey, null, null, null)

    for ((i, r) in result.withIndex()) {
        Log.d("CRYPTO", "result[$i]=${r.toHexString()}")
    }

//    val publicKey = KeyFactory.getInstance("EC").generatePublic(X509EncodedKeySpec(pubKeyBytes))
    // DHKEM_P256_HKDF_SHA256, HKDF_SHA256, HPKE_AES_128_GCM
//    val suiteName = Hpke.getSuiteName(KemParameterSpec.DHKEM_P256_HKDF_SHA256, KdfParameterSpec.HKDF_SHA256,
//        AeadParameterSpec.AES_128_GCM)
//    val suiteName = "DHKEM_P256_HKDF_SHA256/HKDF_SHA256/AES_128_GCM"
//    Log.d("CRYPTO", "HPKE Suite: $suiteName")

//    SignatureConfig.register()
//    val privateKeysetHandle = KeysetHandle.generateNew(
//        EciesParameters.builder()
//            .setVariant(EciesParameters.Variant.NO_PREFIX)
//            .setNistCurvePointFormat(EciesParameters.PointFormat.UNCOMPRESSED)
//            .setCurveType(EciesParameters.CurveType.NIST_P256)
//            .setHashType(EciesParameters.HashType.SHA256)
//            .setDemParameters(
//                AesGcmParameters.builder()
//                    .setIvSizeBytes(12)
//                    .setKeySizeBytes(16)
//                    .setTagSizeBytes(16)
//                    .setVariant(AesGcmParameters.Variant.NO_PREFIX)
//                    .build())
//            .build()
////        HybridKeyTemplates.ECIES_P256_HKDF_HMAC_SHA256_AES128_GCM,
////        Parameters
//    )
//    val stream = ByteArrayOutputStream()
//    val publicKeysetHandle = privateKeysetHandle.publicKeysetHandle
//    CleartextKeysetHandle.write(publicKeysetHandle, JsonKeysetWriter.withOutputStream(stream))
//    Log.d("CRYPTO", "example public key json= $stream")

//    val handle = KeysetHandle.newBuilder()
//        .addEntry(BinaryKeysetReader.withBytes(pubKeyBytes))
//        .build()
//    val hpke = handle.getPrimitive(
//        RegistryConfiguration.get(), HybridEncrypt::class.java)
//    val encrypted = hpke.encrypt(userSecret, info)
//    Log.d("CRYPTO", "encrypted=${encrypted}")

//    val hpke = Hpke.getInstance(suiteName)
//    val message = hpke.seal(publicKey, info, userSecret, null)
//    val kpg = KeyPairGenerator.getInstance("Ed25519", "AndroidKeyStore").generateKeyPair()
//    Log.d("CRYPTO", "Public key: ${kpg.public.encoded}")
//    val kpg = KeyPairGenerator.getInstance(KeyProperties.KEY_ALGORITHM_RSA, "AndroidKeyStore").generateKeyPair()
//    Cipher.getInstance(KeyProperties.KEY_ALGORITHM_RSA)
}