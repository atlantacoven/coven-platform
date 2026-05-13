package space.thecoven.android

import android.content.ComponentName
import android.content.Intent
import android.nfc.NdefMessage
import android.nfc.NfcAdapter
import android.nfc.NfcManager
import android.nfc.cardemulation.CardEmulation
import android.os.Build
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import space.thecoven.android.ui.theme.TheCovenTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        handleAuthIntent(intent)

        val manager = getSystemService(NFC_SERVICE) as NfcManager
        val adapter = manager.defaultAdapter

        Log.d("NFC", "adapter=$adapter enabled=${adapter.isEnabled}")
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            Log.d("NFC", "isSecureNfcSupported=${adapter.isSecureNfcSupported} isSecureNfcEnabled=${adapter.isSecureNfcEnabled}")
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.VANILLA_ICE_CREAM) {
            Log.d("NFC", "isObserveModeSupported=${adapter.isObserveModeSupported} isObserveModeEnabled=${adapter.isObserveModeEnabled} isReaderOptionSupported=${adapter.isReaderOptionSupported} isReaderOptionEnabled=${adapter.isReaderOptionEnabled}")
        }
        if (Build.VERSION.SDK_INT_FULL >= Build.VERSION_CODES_FULL.BAKLAVA_1) {
            Log.d("NFC", "isExitFramesSupported=${adapter.isExitFramesSupported} isTagIntentAllowed=${adapter.isTagIntentAllowed} isTagIntentAppPreferenceSupported=${adapter.isTagIntentAppPreferenceSupported}")
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.UPSIDE_DOWN_CAKE) {
            val antenna = adapter.nfcAntennaInfo
            Log.d("NFC", "w=${antenna?.deviceWidth} h=${antenna?.deviceHeight} foldable=${antenna?.isDeviceFoldable}")
            for (antenna in antenna?.availableNfcAntennas ?: emptyList()) {
                Log.d("NFC", "antenna_x=${antenna.locationX} antenna_y=${antenna.locationY}")
            }
        }

        val cardEm = CardEmulation.getInstance(adapter)

        val isDefault = cardEm.isDefaultServiceForAid(ComponentName(application, DoorAccessService::class.java.name), getString(R.string.aid))

        Log.d("NFC", "isDefault=$isDefault")

//        var message = NdefRecord.createUri("coven://thecoven.space/door")
//        Log.d("NFC", "message=$message id=${message.id.toHexString()} payload-ascii=${message.payload.toString(Charsets.US_ASCII)}")
//        message = NdefRecord.createApplicationRecord(application.packageName)
//        Log.d("NFC", "message=$message id=${message.id.toHexString()} payload-ascii=${message.payload.toString(Charsets.US_ASCII)}")

        enableEdgeToEdge()
        setContent {
            TheCovenTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    Greeting(
                        name = "Android",
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        handleAuthIntent(intent)
    }

    fun handleAuthIntent(intent: Intent) {
        Log.i("NFC", "intent: $intent")
        if (NfcAdapter.ACTION_NDEF_DISCOVERED == intent.action) {
            val messages = intent.getParcelableArrayListExtra<NdefMessage>(NfcAdapter.EXTRA_NDEF_MESSAGES) ?: emptyList<NdefMessage>()
            for (message in messages) {
                Log.i("NFC", "message: $message")
            }
        } else {
            Log.i("NFC", "not an NFC intent")
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello $name!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    TheCovenTheme {
        Greeting("Android")
    }
}
