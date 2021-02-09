#!/usr/bin/env python3

import RPi.GPIO as GPIO
import subprocess, time, socket
from PIL import Image
from Adafruit_Thermal import *


printer = Adafruit_Thermal("/dev/serial0", 19200, timeout=5)

printer.printImage(Image.open('gfx/doron_maze.png'), True)
printer.feed(5)
