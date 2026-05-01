package space.thecoven.android

import android.util.Log
import org.bouncycastle.crypto.hpke.HPKE
import java.security.KeyFactory
import java.security.Signature
import java.security.spec.X509EncodedKeySpec
import kotlin.jvm.Throws


class Authenticator(val doorSigningPubKey: ByteArray) {

    // NOTE: when minSdk is 37 we can use the native Android libs for this instead of bouncycastle
    private val hpke = HPKE(
        HPKE.mode_base,
        HPKE.kem_P256_SHA256,
        HPKE.kdf_HKDF_SHA256,
        HPKE.aead_AES_GCM128
    )

    companion object {
        const val CHALLENGE_SIZE = 8 + 64
        const val NONCE_SIZE = 8

        val INFO_FIELD = "thecoven.space".toByteArray(Charsets.US_ASCII)
    }

    @Throws(SecurityException::class)
    fun verifyChallenge(challenge: ByteArray): ByteArray {
        val nonce = challenge.copyOfRange(0, NONCE_SIZE)
        val signature = challenge.copyOfRange(NONCE_SIZE, challenge.size)

        val keyspec = X509EncodedKeySpec(doorSigningPubKey)
        val key = KeyFactory.getInstance("Ed25519").generatePublic(keyspec)
        val sig = Signature.getInstance("Ed25519")
        sig.initVerify(key)
        sig.update(nonce)
        sig.verify(signature)
        return nonce
    }

    fun authenticate(nonce: ByteArray, userSecret: ByteArray, pubkey: ByteArray): ByteArray {
        val plaintext = nonce + userSecret

        val pk = hpke.deserializePublicKey(pubkey)
        val aar = byteArrayOf()
        val result = hpke.seal(pk, INFO_FIELD, aar, plaintext, null, null, null)

        val cipherText = result[0]
        val encapsulatedKey = result[1]
        return cipherText+encapsulatedKey
    }
}
