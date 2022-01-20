#!/usr/bin/env python3
# Must be ran with sudo! 

from prometheus_client import make_wsgi_app, Gauge
from wsgiref.simple_server import make_server
from sys import argv
import glob, os, subprocess

def main():
    # Create Gauge
    g = Gauge('sauna_temperature', 'Check temperature')

    os.system('modprobe w1-gpio')
    os.system('modprobe w1-therm')

    base_dir = '/sys/bus/w1/devices/'
    device_folder = glob.glob(base_dir + '28*')[0]
    device_file = device_folder + '/w1_slave'

    def read_temp_raw():
        f = open(device_file, 'r')
        lines = f.readlines()
        f.close()
        return lines

    def read_temp():
        lines = read_temp_raw()
        while lines[0].strip()[-3:] != 'YES':
            time.sleep(0.2)
            lines = read_temp_raw()
        equals_pos = lines[1].find('t=')
        if equals_pos != -1:
            temp_string = lines[1][equals_pos+2:]
            temp = float(temp_string) / 1000.0
            return temp

    metrics_app = make_wsgi_app()

    def my_app(environ, start_fn):
        if environ['PATH_INFO'] == '/metrics':
            g.set(read_temp())
            return metrics_app(environ, start_fn)

    httpd = make_server('', 8000, my_app)
    httpd.serve_forever()

if __name__ == "__main__":
    try:
        argv[1]
        print(f"usage: sudo python3 {argv[0]}\n\nExporter for Prometheus made for DS18B20 sensor.")

    except IndexError:
        main()