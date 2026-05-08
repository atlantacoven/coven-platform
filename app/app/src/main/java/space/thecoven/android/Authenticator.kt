package space.thecoven.android

import android.util.Log
import org.bouncycastle.crypto.hpke.HPKE
import org.bouncycastle.math.ec.rfc8032.Ed25519
import java.security.KeyFactory
import java.security.Signature
import java.security.spec.X509EncodedKeySpec
import kotlin.jvm.Throws


class Authenticator(val doorSigningPubKey: ByteArray) {

    // NOTE: when minSdk is 37 we can use the native Android libs for this instead of bouncycastle
    private val hpke = HPKE(
        HPKE.mode_base,
        HPKE.kem_X25519_SHA256,
        HPKE.kdf_HKDF_SHA256,
        HPKE.aead_AES_GCM128
    )

    companion object {
        const val UINT64_SIZE = 8
        const val CHALLENGE_SIZE = UINT64_SIZE + Ed25519.SIGNATURE_SIZE
        const val NONCE_SIZE = UINT64_SIZE

        val INFO_FIELD = "thecoven.space".encodeToByteArray()
    }

    sealed class ChallengeResult {
        object InvalidSignature : ChallengeResult()
        class ValidSignature(val nonce: ByteArray) : ChallengeResult()
    }

    @Throws(SecurityException::class)
    fun verifyChallenge(challenge: ByteArray): ChallengeResult {
        val nonce = challenge.copyOfRange(0, NONCE_SIZE)
        val signature = challenge.copyOfRange(NONCE_SIZE, challenge.size)

        val keyspec = X509EncodedKeySpec(doorSigningPubKey)
        val key = KeyFactory.getInstance("Ed25519").generatePublic(keyspec)
        val sig = Signature.getInstance("Ed25519")
        sig.initVerify(key)
        sig.update(nonce)
        val valid = sig.verify(signature)
        if (valid) {
            return ChallengeResult.ValidSignature(nonce)
        } else {
            return ChallengeResult.InvalidSignature
        }
    }

    fun authenticate(nonce: ByteArray, userSecret: ByteArray, pubkey: ByteArray): ByteArray {
        val plaintext = nonce + userSecret

//        val skS = hpke.generatePrivateKey()
//        val senderPublicKey = hpke.serializePublicKey(skS.public)

        val pkR = hpke.deserializePublicKey(pubkey)
        val aar = byteArrayOf() // empty
        val result = hpke.seal(pkR,
            INFO_FIELD, aar,
            plaintext,
            null, null, // PSK (unused)
            null)

        val cipherText = result[0]
        val encapsulatedKey = result[1]
        Log.d("CRYPTO", "cipher=${cipherText.toHexString()}")
        Log.d("CRYPTO", "encapsulatedKey=${encapsulatedKey.toHexString()}")
//        Log.d("CRYPTO", "pks=${senderPublicKey.toHexString()}")
        return cipherText+encapsulatedKey
    }
}
