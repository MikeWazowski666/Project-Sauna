#!/bin/bash

# Grafana setup
wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | tee -a /etc/apt/sources.list.d/grafana.list
apt-get install -y grafana

# Start Grafana
/bin/systemctl enable grafana-server
/bin/systemctl start grafana-server

# Do Prometheus setup
wget https://github.com/prometheus/prometheus/releases/download/v2.22.0/prometheus-2.22.0.linux-armv7.tar.gz
tar xfz prometheus-2.22.0.linux-armv7.tar.gz
mv prometheus-2.22.0.linux-armv7/ ~/prometheus/
rm prometheus-2.22.0.linux-armv7.tar.gz
touch /etc/systemctl/system/prometheus.service
echo "
[Unit]
Description=Prometheus Server
Documentation=https://prometheus.io/docs/introduction/overview/
After=network-online.target

[Service]
User=pi
Restart=on-failure

ExecStart=/home/pi/prometheus/prometheus \
  --config.file=/home/pi/prometheus/prometheus.yml \
  --storage.tsdb.path=/home/pi/prometheus/data

[Install]
WantedBy=multi-user.target
" > /etc/systemctl/system/prometheus.service

# Start prometheus
/bin/systemctl enable prometheus
/bin/systemctl start prometheus

# Install Python dependancies
apt install -y python3 python3-pip
python3 -m pip install -r requirements.txt

# Add crontab to mantain data
echo "@reboot /usr/bin/python3 /home/pi/pi/run.py" >> /var/spool/crontab/crontabs/root
