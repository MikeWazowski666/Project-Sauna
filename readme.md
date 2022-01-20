# Project Sauna

# Description
My family wanted to get more information about the temperatures in sauna. 
So I created a digital temperature sensor, that saves the temperature data in [Grafana](https://grafana.com/) with [Promertheus](https://prometheus.io/) and a custom exporter.

# What is needed for this project?

- Raspberry Pi with Internet connection
- SD card
- DS18B20 (Waterproof Temperature sensor)
- 1 x 4.7k Ohm resistor
- Jumper wires

# How to wire it up?
Check out [this](Schematic.pdf) schematic.

# Installation

- Install NOOBS or any other OS on your SD card. Make sure it is Raspberry Pi compatible!
- Clone this repo
- Run `sudo setup.sh` (currently in development)
- Finally run `sudo pi/run.py` (in development)
