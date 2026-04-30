package space.thecoven.android

import android.security.keystore.KeyProperties
import java.security.KeyPairGenerator
import java.security.KeyStore
import javax.crypto.Cipher

// Server needs to give App a user-unique secret.
//    Server->App: signature=sign(user_id+salt, server_private_key)
//          user_secret=user_id+salt+signature
//          salt must be at least 40 bits
// App needs to verify it's talking to the Door. Door needs to verify the App generated the key and did so recently.
//    Door->App: door_public_key+nonce
//    App->Door: access_key=encrypt(user_id+nonce+user_secret, door_public_key)
//    Door: nonce, user_secret=decrypt(access_key, door_private_key)
//      user_id,salt,signature=user_secret
// Door needs to verify the Server gave out the key.
//      verify(signature,user_id+salt, server_pub_key)

fun generateKey() {
    val ks = KeyStore.getInstance("AndroidKeyStore").also { it.load(null) }

    val userKey = "0000000000001234bd0acb5ff2d5f065fb564f0d57188c2ce1151ea6a6a8f4c069f3474b647745158d5c64aee927c8c277f075f5225f8ef3247415b1e72018e9fcc6364660329f8acf5a5ed30fd984e4a994de39ebbbecaaa65d9cb27b100b0d".hexToByteArray()

    ks.setKeyEntry("UserSecret", userKey, emptyArray())

//    val kpg = KeyPairGenerator.getInstance("Ed25519", "AndroidKeyStore").generateKeyPair()
//    Log.d("CRYPTO", "Public key: ${kpg.public.encoded}")

    val kpg = KeyPairGenerator.getInstance(KeyProperties.KEY_ALGORITHM_RSA, "AndroidKeyStore").generateKeyPair()

    Cipher.getInstance(KeyProperties.KEY_ALGORITHM_RSA)
}