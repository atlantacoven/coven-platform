package space.thecoven.android

import android.nfc.cardemulation.HostApduService
import android.os.Bundle
import android.util.Log
import kotlin.collections.toByteArray

class DoorAccessService : HostApduService() {

    override fun onCreate() {
        super.onCreate()
        Log.d("NFC", "DoorAccessService started")
    }

    override fun processCommandApdu(commandApdu: ByteArray, extras: Bundle?): ByteArray {
        Log.d("NFC", "processCommandApdu $commandApdu extras=$extras")
        // this command blocks the main thread. If it can't be executed immediately,
        // return null and call sendResponseApdu() when ready
        return arrayOf(0xDE, 0xAD, 0xBE, 0xEF).map(Int::toByte).toByteArray()
    }

    override fun onDeactivated(reason: Int) {
        Log.w("NFC", "DoorAccessService deactivated reason=${reason}")
    }
}
