#!/usr/bin/env python3
#
#
#
# --------------------
# Must be ran as root! 
#---------------------
#
#
#

from prometheus_client import make_wsgi_app, Gauge
from wsgiref.simple_server import make_server
import glob, os

def main():
    # Create Gauge
    g = Gauge('sauna_temperature', 'Check temperature', ['sensor'])

    # Make slaves
    os.system('modprobe w1-gpio')
    os.system('modprobe w1-therm')

    dev_name = []
    dev_file = []
    
    # Search for devices
    base_dir = '/sys/bus/w1/devices/'
    dev_folder = glob.glob(base_dir + '28*')
    for i in range(len(dev_folder)):
        dev_name.append(dev_folder[i][23:])
        dev_file.append(dev_folder[i] + '/temperature')

    # Read temperature    
    def read_temp(fn):
        with open(fn) as h:
            temp = float(h.read()) / 1000 # Convert to celsus
            return temp
    
    def web_metrics(environ, start_fn):
        # Check if path is to `/metrics`
        if environ['PATH_INFO'] == '/metrics':
            # Set temperature to lable
            for i in range(len(dev_folder)):
                g.labels(dev_name[i]).set(read_temp(dev_file[i]))
            # Return
            return metrics_app(environ, start_fn)
    
    # Run exporter
    metrics_app = make_wsgi_app()
    httpd = make_server('', 8000, web_metrics)
    httpd.serve_forever()

if __name__ == "__main__":
    main()
