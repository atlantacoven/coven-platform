package space.thecoven.android

import org.bouncycastle.crypto.hpke.HPKE
import org.junit.Test

import org.junit.Assert.*

/**
 * Example local unit test, which will execute on the development machine (host).
 *
 * See [testing documentation](http://d.android.com/tools/testing).
 */
class ExampleUnitTest {
    @Test
    fun addition_isCorrect() {
        assertEquals(4, 2 + 2)
    }

    @Test
    fun hpkeAuth() {
        val hpke = HPKE(
            HPKE.mode_base,
            HPKE.kem_X25519_SHA256,
            HPKE.kdf_HKDF_SHA256,
            HPKE.aead_AES_GCM128
        )

        val info = "thecoven.space".encodeToByteArray() // byteArrayOf(0xDE, 0xAD, 0xBE, 0xEF)
        val aar = byteArrayOf()

        // receiver
        val skR = hpke.generatePrivateKey()
        val pkR = hpke.serializePublicKey(skR.public)

        // sender
        val messageIn = "shh this is secret"
        val pt = messageIn.encodeToByteArray()
        val (ct, enc) = hpke.seal(hpke.deserializePublicKey(pkR), info, aar, pt, null, null, null)

        // receiver
        val res = hpke.open(enc, skR, info, aar, ct, null, null, null)

        val messageOut = res.decodeToString()
        assertEquals(messageOut, messageIn)
    }
}