#!/usr/bin/env python3

# Weather forecast for Raspberry Pi w/Adafruit Mini Thermal Printer.
# Retrieves data from DarkSky.net's API, prints current conditions and
# forecasts for next two days.  See timetemp.py for a different
# weather example using nice bitmaps.
# Written by Adafruit Industries.  MIT license.
# 
# Required software includes Adafruit_Thermal and PySerial libraries.
# Other libraries used are part of stock Python install.
# 
# Resources:
# http://www.adafruit.com/products/597 Mini Thermal Receipt Printer
# http://www.adafruit.com/products/600 Printer starter pack

from random import choice
from Adafruit_Thermal import *
import urllib, json

printer = Adafruit_Thermal("/dev/serial0", 19200, timeout=5)

url = "https://niceties.herokuapp.com"
response = urllib.urlopen(url)
data = json.loads(response.read())
recursors = [
    "Bobby (Robert) DeLanghe"
]

# Print nice thing
printer.inverseOn()
printer.print(' ' + choice(recursors) + ' \n')
printer.inverseOff()
printer.print(data)

# Print feed
printer.feed(6)
