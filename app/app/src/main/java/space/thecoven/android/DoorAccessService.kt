package space.thecoven.android

import android.nfc.cardemulation.HostApduService
import android.os.Bundle
import android.util.Log

/**
 * This class uses Android's support for Host Card Emulation to exchange keys with the door lock
 * even from the background. Communication is handled through standards ISO-14443-4 (data-link)
 * and ISO-7816-4 (application).
 *
 * NOTE: the implementation goes well above and beyond what is strictly necessary to implement
 * the protocol. I've done so intentionally do document the ISO-7816 protocol.
 *
 * See: https://developer.android.com/develop/connectivity/nfc/hce
 * See: https://www.freecalypso.org/pub/GSM/ISO7816/ISO_7816-4_2005.pdf
 */
class DoorAccessService : HostApduService() {

    companion object {
        const val DOOR_UNLOCK_RESULT_CMD = 0xFA
        val USER_SECRET = "0000000000001234aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa609510891e090572c105d0ce3dccf8125e4403daea4284de511f10b4d154b1a8f48c3a8a9cd99aca946561aedf4a6eb562741f140687835276e6f698ea499407".hexToByteArray() // STOPSHIP
    }

    override fun onCreate() {
        super.onCreate()
        Log.d("NFC", "DoorAccessService started")
    }

    override fun processCommandApdu(commandApdu: ByteArray, extras: Bundle?): ByteArray {
        Log.d("NFC", "processCommandApdu ${commandApdu.toHexString()} extras=$extras")
        // this command blocks the main thread. If it can't be executed immediately,
        // return null and call sendResponseApdu() when ready
        val res = try {
            val apdu = IDCard.CommandAPDU.decode(commandApdu)
            when (apdu) {
                is IDCard.CommandAPDU.Proprietary ->
                    Log.d("NFC", "RequestADPU ${apdu.cla.toUByte()} (proprietary) ${apdu.raw.toHexString()}")

                is IDCard.CommandAPDU.InterIndustry ->
                    Log.d(
                        "NFC",
                        "RequestADPU ${apdu.command.name} ${apdu.data.toString(Charsets.US_ASCII)}"
                    )
            }
            handleCommand(apdu)
        } catch (e: IDCard.InvalidMessageException) {
            Log.e("NFC", "Invalid message", e)
            IDCard.Status.InvalidLength.toResponse()
        }
        Log.d("NFC", "ResponseADPU status=${res.status} || ${res.encode().toHexString()}")
        return res.encode()
    }

    fun handleCommand(cmd: IDCard.CommandAPDU): IDCard.ResponseAPDU {
        when (cmd) {
            is IDCard.CommandAPDU.InterIndustry -> {
                if (cmd.channel != IDCard.LogicalChannel.ZERO)
                    return IDCard.Status.UnsupportedCLAFunctionLogicalChannel.toResponse()
                if (cmd.chaining.toInt() != 0)
                    return IDCard.Status.UnsupportedCLAFunctionCommandChanning.toResponse()
                if (cmd.secureMessaging.toInt() != 0)
                    return IDCard.Status.UnsupportedCLAFunctionSecureMessaging.toResponse()

                when (cmd.command) {
                    IDCard.Command.SELECT -> {
                        val aid = cmd.data.toHexString(HexFormat.UpperCase)
                        if (aid != getString(R.string.aid))
                            return IDCard.Status.FileOrApplicationNotFound.toResponse()
                        return IDCard.Status.OK.toResponse()
                    }

                    IDCard.Command.GENERAL_AUTHENTICATE -> {
                        val challenge = cmd.data.copyOfRange(0, Authenticator.CHALLENGE_SIZE)
                        val pubKey = cmd.data.copyOfRange(Authenticator.CHALLENGE_SIZE, cmd.data.size)
                        Log.d("NFC", "challenge=${challenge.toHexString()}")
                        Log.d("NFC", "pubKey=${pubKey.toHexString()}")

                        val doorPubSigningKey = getString(R.string.doorPubSigningKey).hexToByteArray()
                        val auth = Authenticator(doorPubSigningKey)

                        try {
                            val nonce = auth.verifyChallenge(challenge)
                            Log.d("NFC", "nonce=${nonce.toHexString()}")
                            val accessKey = auth.authenticate(nonce, USER_SECRET, pubKey)
                            Log.d("NFC", "accessKey=${accessKey.toHexString()}")

                            return IDCard.ResponseAPDU(data = accessKey)
                        } catch (e: Exception) {
                            Log.e("NFC", "Authentication Failed", e)
                            // Return a security error
                            return IDCard.Status(0x66.toUByte(), 0x00.toUByte()).toResponse()
                        }
                    }

                    else ->
                        return IDCard.Status.UnsupportedCommand.toResponse()
                }
            }

            is IDCard.CommandAPDU.Proprietary -> {
                when (val proprietaryCommand = cmd.cla.toUByte().toInt()) {
                    DOOR_UNLOCK_RESULT_CMD -> {
                        // auth result
                        if (cmd.raw.size < 2) {
                            Log.d("NFC", "invalid command size=${cmd.raw.size}")
                            return IDCard.Status.InvalidLength.toResponse()
                        }
                        val status = IDCard.Status(cmd.raw[0].toUByte(), cmd.raw[1].toUByte())
                        if (status.isOkay()) {
                            Log.d("NFC", "door unlocked")
                            // TODO: show local notification
                        } else {
                            Log.d("NFC", "door NOT unlocked: $status")
                        }
                        return IDCard.Status.OK.toResponse()
                    }
                    0xFF -> {
                        // Dummy test command
                        return IDCard.ResponseAPDU(data = listOf(0xDE, 0xAD, 0xBE, 0xEF).map { it.toByte() }.toByteArray())
                    }

                    else -> {
                        Log.d("NFC", "unknown command $proprietaryCommand")
                        return IDCard.Status.UnsupportedCommand.toResponse()
                    }
                }
            }
        }
    }

    override fun onDeactivated(reason: Int) {
        when (reason) {
            DEACTIVATION_LINK_LOSS ->
                Log.w("NFC", "DoorAccessService deactivated reason=link lost")

            DEACTIVATION_DESELECTED ->
                Log.w("NFC", "DoorAccessService deactivated reason=deactivated")

            else ->
                Log.w("NFC", "DoorAccessService deactivated reason=unknown ($reason)")
        }
    }
}
