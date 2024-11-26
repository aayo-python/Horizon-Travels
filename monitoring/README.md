# Higher Learning Startup - Monitoring(ft/monitoring)
<!-- https://betterstack.com/community/guides/monitoring/monitor-linux-prometheus-node-exporter/#step-4-configuring-prometheus -->

## This covers the process of installation, configuration and alerting on EC2 Server

### Install prometheus on windows
```
prometheus.io/download
unzip
cmd // cd D:/unzipped_folder
prometheus.exe
```
<!-- server should be ready on localhost:9090 -->

### Collecting metrics on linux server
```
ssh -i "prometheus.pem" ubuntu@ec2-xyz.southeast...
sudo apt-get update
sudo apt-get upgrade
sudo wget https://github.com/prometheus/node_exporter/releases/download/v1.8.2/node_exporter-1.8.2.darwin-amd64.tar.gz
sudo tar -xvf node_exporter-1.8.2.darwin-amd64.tar.gz
cd node_exporter
./node_exporter
```
<!-- should be listening on localhost:9100 -->
<!-- now make sure this port is opened to windows machine sudo ufw allow 9100 -->

### Add the linux server to main prometheus server
```
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]

```