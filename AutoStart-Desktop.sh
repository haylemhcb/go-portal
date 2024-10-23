#!/bin/bash

pkill portalrobot
pkill dhclient
pkill airodump

x-terminal-emulator -e './run-all.sh'&
sleep 30
x-www-browser http://localhost:9742


pkill portalrobot
pkill dhclient
pkill airodump
