import rp2
from machine import Pin
import ubinascii
from utime import sleep
import network
import json

led_p = Pin("LED", Pin.OUT)
relay_p = Pin(8, Pin.OUT)
beeper_p = Pin(9, Pin.OUT)


def connect_wifi():
    rp2.country("US")
    mac = ubinascii.hexlify(network.WLAN().config('mac'),':').decode()
    print("mac = " + mac)

    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    with open('wifi.json', 'r') as f:
        config = json.load(f)
    wlan.connect(config['ssid'], config['password'])

    # wait for WIFI connect
    max_wait = 30
    while max_wait > 0:
        if wlan.status() < 0 or wlan.status() >= 3:
            break
        max_wait -= 1
        print('waiting for connection...')
        sleep(1)

    # Handle connection error
    if wlan.status() != 3:
        raise RuntimeError('network connection failed')
    else:
        print('connected')
        status = wlan.ifconfig()
        print( 'ip = ' + status[0] )

def unlock_door():
    relay_p.low()

def lock_door():
    relay_p.high()

def unlock_temporarily(timeout = 5):
    unlock_door()
    sleep(timeout)
    lock_door()

def beep_on():
    beeper_p.on()

def beep_off():
    beeper_p.off()

def main():
    # connect_wifi()
    # lock_door()
    while True:
        lock_door()
        led_p.toggle()
        beep_on()
        sleep(5)
        unlock_door()
        led_p.toggle()
        beep_off()
        sleep(5)

main()
