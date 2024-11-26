## Step 1 — Preparing the servers
Before proceeding with the installation of Prometheus and Node Exporter, it is necessary to prepare your servers by creating users and directories, opening ports in the firewall, and generating a certificate that will secure the connection between the Prometheus server and the Node Exporter server.


### Server 1
Log in to the server and create a user named prometheus through the commands below:

<!-- ssh <your_first_server_ip> -->
```
ssh -i "prometheus.pem" ubuntu@ec2-xyz.southeast...
sudo apt-get update
sudo apt-get upgrade
sudo useradd --no-create-home --shell /bin/false prometheus
```
Here, the options --no-create-home and --shell /bin/false will ensure that no home directory is created for this user, and this user cannot login to the server.

Next, create the following directories for storing Prometheus' configuration files and data: 

```
sudo mkdir /etc/prometheus
```

Then, change the user and group ownership on the new directories to the prometheus user.
```
sudo mkdir /var/lib/
sudo chown -R prometheus:prometheus /etc/prometheus /var/lib/prometheus
 ```

If you plan to run Prometheus on the local network and do not use a firewall, you can skip these

Open the default HTTP port for the Prometheus server:
```
sudo ufw allow 9090/tcp
sudo ufw reload
```

At this point, we are done with the initial preparation of the Prometheus server. 



### Server 2

Let's go ahead and prepare the second system, which will run a Node Exporter instance. Login to the server and create a new user called nodeuser:

<!-- ssh <your_second_server_ip> -->
```
ssh -i "prometheus.pem" ubuntu@ec2-xyz.southeast...
sudo apt-get update
sudo apt-get upgrade
sudo useradd --no-create-home --shell /bin/false nodeuser
``` 

Next, create a directory for storing the Node Exporter's configuration files:

```
sudo mkdir /etc/node_exporter
```

Node Exporter uses port 9100 for incoming connections, so the Prometheus server will connect to this port. If you use a firewall, run the commands below on the Node Exporter server to allow connections on this port:

```
sudo ufw allow 9100/tcp
sudo ufw reload
```



## Step 2 — Securing the connection between the servers
We need to create an SSL certificate to secure the connections between the Prometheus server and the Node Exporter. Here, we will be using self-signed credentials. Run the following command on the Node Exporter server to create a certificate:

```
openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout node_exporter.key -out node_exporter.crt -subj "/C=US/ST=Utah/L=Lehi/O=Your Company, Inc./OU=IT/CN=yourdomain.com" -addext "subjectAltName = IP:<your_second_server_ip>"
```

Don't forget to replace <your_second_server_ip> string with the IP address of your second server, to which the Prometheus server will connect. You will see the following output:

```
 Output
Generating a RSA private key
...+++++
...............................................+++++
writing new private key to 'node_exporter.key'
-----
```

As a result, you will get two files in the current directory: 
    node_exporter.crt 
    node_exporter.key. 
Copy these files to the /etc/node_exporter directory:
 
```
sudo cp node_exporter.*  /etc/node_exporter
```
Afterward, change the user and group ownership on the new directories and files to nodeuser.

```
sudo chown -R nodeuser:nodeuser /etc/node_exporter
```

Then copy the node_exporter.crt file to the first server in any way you can. For example, you can show the contents of a file and copy it to the clipboard:

```
cat /etc/node_exporter/node_exporter.crt
```
 Output
 ```
-----BEGIN CERTIFICATE-----
MIIDzjCCAragAwIBAgIUCQuU1Y/15mJiRHF5zKc5nPql8wUwDQYJKoZIhvcNAQEL
BQAwbjELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFV0YWgxDTALBgNVBAcMBExlaGkx
GzAZBgNVBAoMEllvdXIgQ29tcGFueSwgSW5jLjELMAkGA1UECwwCSVQxFzAVBgNV
BAMMDnlvdXJkb21haW4uY29tMB4XDTIyMDQxMzE3MjIyM1oXDTIzMDQxMzE3MjIy
M1owbjELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFV0YWgxDTALBgNVBAcMBExlaGkx
GzAZBgNVBAoMEllvdXIgQ29tcGFueSwgSW5jLjELMAkGA1UECwwCSVQxFzAVBgNV
BAMMDnlvdXJkb21haW4uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqXopi/6gWVx7qe4ADPc1g2VK8uEcCVLj6EPlfTmatxK8AlUGMFx8zP5H8xrb
7S3sdUN+iGQAYTLavoUIBR2k8IBmq073ziuK+Dd56TYvcjdJvOXyKgRo2R49NdtN
3uYtx3ZZp9E1QmMQMuMev4qCQ0jIwVtxb0jEk8RimAQFoFc5z86JD9oFSVp+VqA3
HPUElmr5at7+Nczf07m874S47QfvhHIiUDl8PX/Ch+gm6gUBqic7adltz+pgkl1W
KyvBGGbK48QBR57rMn0Juy3DHHgQyEYNXWUlQ2gjc1UuK0NZZW8KTXY29jF8mwZh
tdZ/UFiqoIgMfRry3522AnAZwQIDAQABo2QwYjAdBgNVHQ4EFgQUF3juV9xC32zU
S33qwuo1V6LIPNowHwYDVR0jBBgwFoAUF3juV9xC32zUS33qwuo1V6LIPNowDwYD
VR0TAQH/BAUwAwEB/zAPBgNVHREECDAGhwSHtW71MA0GCSqGSIb3DQEBCwUAA4IB
AQARRTXVPzgpKWfgkF4dG0Cqd/HlmdRGFGU4S7LGWdbHSW6B7s0A4YvXOzQmL+ib
ucnbbsJn4Vy4GtOkpkfJgzSQqqm0ozcZsKOJdRkZ7/S8NNdVI/leqO5QT9I0kinQ
NfTd544f/i4MCBhvGx0yhQum32xIJy75RcQC4m8y76dzlW8OCyFidcw+S9plcsml
JIOV0tzfeYc8vQrHMGyhOrJk9FCCtnipDd5AtlFqB1Q3J2bz+7A4JvG2uK2qGOAI
JuiBuqDpNGgIotV3LcqzkAmIT0JomTDx9MGCfIJFz71iO80JXGHkk/L4U2U+9cvQ
Gllx69QjPJirRJj2LNyEB0PD
-----END CERTIFICATE-----
```

Then create a new file on the Prometheus server and paste the contents there:

```
sudo nano /etc/prometheus/node_exporter.crt

-----BEGIN CERTIFICATE-----
MIIDzjCCAragAwIBAgIUCQuU1Y/15mJiRHF5zKc5nPql8wUwDQYJKoZIhvcNAQEL
BQAwbjELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFV0YWgxDTALBgNVBAcMBExlaGkx
GzAZBgNVBAoMEllvdXIgQ29tcGFueSwgSW5jLjELMAkGA1UECwwCSVQxFzAVBgNV
BAMMDnlvdXJkb21haW4uY29tMB4XDTIyMDQxMzE3MjIyM1oXDTIzMDQxMzE3MjIy
M1owbjELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFV0YWgxDTALBgNVBAcMBExlaGkx
GzAZBgNVBAoMEllvdXIgQ29tcGFueSwgSW5jLjELMAkGA1UECwwCSVQxFzAVBgNV
BAMMDnlvdXJkb21haW4uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqXopi/6gWVx7qe4ADPc1g2VK8uEcCVLj6EPlfTmatxK8AlUGMFx8zP5H8xrb
7S3sdUN+iGQAYTLavoUIBR2k8IBmq073ziuK+Dd56TYvcjdJvOXyKgRo2R49NdtN
3uYtx3ZZp9E1QmMQMuMev4qCQ0jIwVtxb0jEk8RimAQFoFc5z86JD9oFSVp+VqA3
HPUElmr5at7+Nczf07m874S47QfvhHIiUDl8PX/Ch+gm6gUBqic7adltz+pgkl1W
KyvBGGbK48QBR57rMn0Juy3DHHgQyEYNXWUlQ2gjc1UuK0NZZW8KTXY29jF8mwZh
tdZ/UFiqoIgMfRry3522AnAZwQIDAQABo2QwYjAdBgNVHQ4EFgQUF3juV9xC32zU
S33qwuo1V6LIPNowHwYDVR0jBBgwFoAUF3juV9xC32zUS33qwuo1V6LIPNowDwYD
VR0TAQH/BAUwAwEB/zAPBgNVHREECDAGhwSHtW71MA0GCSqGSIb3DQEBCwUAA4IB
AQARRTXVPzgpKWfgkF4dG0Cqd/HlmdRGFGU4S7LGWdbHSW6B7s0A4YvXOzQmL+ib
ucnbbsJn4Vy4GtOkpkfJgzSQqqm0ozcZsKOJdRkZ7/S8NNdVI/leqO5QT9I0kinQ
NfTd544f/i4MCBhvGx0yhQum32xIJy75RcQC4m8y76dzlW8OCyFidcw+S9plcsml
JIOV0tzfeYc8vQrHMGyhOrJk9FCCtnipDd5AtlFqB1Q3J2bz+7A4JvG2uK2qGOAI
JuiBuqDpNGgIotV3LcqzkAmIT0JomTDx9MGCfIJFz71iO80JXGHkk/L4U2U+9cvQ
Gllx69QjPJirRJj2LNyEB0PD
-----END CERTIFICATE-----
```

Also, on both servers, you will need the htpasswd utility, which is part of the apache2-utils package.

```
sudo apt install apache2-utils
```

This concludes the preliminary preparation and we can proceed to installing Prometheus on the first server in the next step.



### Step 3 — Downloading and installing Prometheus
Prometheus cannot be installed from the default Ubuntu repositories at the time of writing, so you need to get an archive of the current version. Go to the official download page  and copy the link to the Linux package of the latest stable release.

Prometheus download link 
```
prometheus.io/download
```
Then log in to your first Linux server and download the archive with wget:

```
wget https://github.com/prometheus/prometheus/releases/download/v2.34.0/prometheus-2.34.0.linux-amd64.tar.gz
```

Once the archive is saved, we will verify the integrity of the download by running the sha256sum program and compare its output with the corresponding SHA256 Checksum value listed on the site:

```
sha256sum prometheus-2.34.0.linux-amd64.tar.gz
```
The output should look like this:

Output

```
9ec560940bf53361dd9d3a867d51ceb96f3854ae12f5e532b7d3f60c27f364d0 prometheus-2.34.0.linux-amd64.tar.gz
```

Once you've confirmed that the value above matches the SHA256 Checksum on the website, you can go ahead and unpack the archive using the command shown below:

```
tar zxvf prometheus-2.34.0.linux-amd64.tar.gz
``` 

Output
```
prometheus-2.34.0.linux-amd64/
prometheus-2.34.0.linux-amd64/consoles/
prometheus-2.34.0.linux-amd64/consoles/index.html.example
prometheus-2.34.0.linux-amd64/consoles/node-cpu.html
prometheus-2.34.0.linux-amd64/consoles/node-disk.html
prometheus-2.34.0.linux-amd64/consoles/node-overview.html
prometheus-2.34.0.linux-amd64/consoles/node.html
prometheus-2.34.0.linux-amd64/consoles/prometheus-overview.html
prometheus-2.34.0.linux-amd64/consoles/prometheus.html
prometheus-2.34.0.linux-amd64/console_libraries/
prometheus-2.34.0.linux-amd64/console_libraries/menu.lib
prometheus-2.34.0.linux-amd64/console_libraries/prom.lib
prometheus-2.34.0.linux-amd64/prometheus.yml
prometheus-2.34.0.linux-amd64/LICENSE
prometheus-2.34.0.linux-amd64/NOTICE
prometheus-2.34.0.linux-amd64/prometheus
prometheus-2.34.0.linux-amd64/promtool
```

Go ahead and change into the newly created directory:

```
cd prometheus-2.34.0.linux-amd64/
```

Let's distribute the files in this directory to their appropriate location. Copy the two binary files in the current directory to /usr/local/bin:

 ```
sudo cp prometheus promtool /usr/local/bin/
```

Copy the consoles and console_libraries directories to /etc/prometheus. Also, copy the sample configuration file there:

 
```
sudo cp -r console_libraries consoles prometheus.yml /etc/prometheus
```

Now change the owner for all copied files to prometheus:
 
```
sudo chown prometheus:prometheus /usr/local/bin/prometheus /usr/local/bin/promtool
sudo chown -R prometheus:prometheus /etc/prometheus
``` 

Prometheus is now installed on your server! Before starting it, let's take a look at its configuration file in the next step.



### Step 4 — Configuring Prometheus
In the previous step, we copied the sample config file to the /etc/prometheus directory. Go ahead and open it in your text editor for examination.

```
sudo nano /etc/prometheus/prometheus.yml
```
You should observe the following contents:

/etc/prometheus/prometheus.yml
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


### Step 5 — Running the Prometheus CLI
First, we will start Prometheus manually, and then we will create a Systemd service to manage it properly.

Run the following command on your first server:

```
sudo -u prometheus /usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries
```

You should observe the following launch log:

Output

```ts=2022-04-14T07:35:54.363Z caller=main.go:479 level=info msg="No time or size retention was set so using the default time retention" duration=15d
ts=2022-04-14T07:35:54.364Z caller=main.go:516 level=info msg="Starting Prometheus" version="(version=2.34.0, branch=HEAD, revision=881111fec4332c33094a6fb2680c71fffc427275)"
ts=2022-04-14T07:35:54.364Z caller=main.go:521 level=info build_context="(go=go1.17.8, user=root@121ad7ea5487, date=20220315-15:18:00)"
ts=2022-04-14T07:35:54.364Z caller=main.go:522 level=info host_details="(Linux 5.4.0-104-generic #118-Ubuntu SMP Wed Mar 2 19:02:41 UTC 2022 x86_64 ubuntu-2gb-nbg1-1 (none))"
ts=2022-04-14T07:35:54.365Z caller=main.go:523 level=info fd_limits="(soft=1024, hard=1048576)"
ts=2022-04-14T07:35:54.365Z caller=main.go:524 level=info vm_limits="(soft=unlimited, hard=unlimited)"
ts=2022-04-14T07:35:54.378Z caller=web.go:540 level=info component=web msg="Start listening for connections" address=0.0.0.0:9090
ts=2022-04-14T07:35:54.387Z caller=main.go:937 level=info msg="Starting TSDB ..."
ts=2022-04-14T07:35:54.394Z caller=tls_config.go:195 level=info component=web msg="TLS is disabled." http2=false
ts=2022-04-14T07:35:54.405Z caller=head.go:493 level=info component=tsdb msg="Replaying on-disk memory mappable chunks if any"
ts=2022-04-14T07:35:54.406Z caller=head.go:536 level=info component=tsdb msg="On-disk memory mappable chunks replay completed" duration=4.385µs
ts=2022-04-14T07:35:54.406Z caller=head.go:542 level=info component=tsdb msg="Replaying WAL, this may take a while"
ts=2022-04-14T07:35:54.407Z caller=head.go:613 level=info component=tsdb msg="WAL segment loaded" segment=0 maxSegment=0
ts=2022-04-14T07:35:54.407Z caller=head.go:619 level=info component=tsdb msg="WAL replay completed" checkpoint_replay_duration=49.735µs wal_replay_duration=1.364197ms total_replay_duration=1.547724ms
ts=2022-04-14T07:35:54.409Z caller=main.go:958 level=info fs_type=EXT4_SUPER_MAGIC
ts=2022-04-14T07:35:54.409Z caller=main.go:961 level=info msg="TSDB started"
ts=2022-04-14T07:35:54.409Z caller=main.go:1142 level=info msg="Loading configuration file" filename=/etc/prometheus/prometheus.yml
ts=2022-04-14T07:35:54.413Z caller=main.go:1179 level=info msg="Completed loading of configuration file" filename=/etc/prometheus/prometheus.yml totalDuration=3.377506ms db_storage=1.585µs remote_storage=2.659µs web_handler=718ns query_engine=8.708µs scrape=2.301903ms scrape_sd=31.007µs notify=386.697µs notify_sd=21.135µs rules=1.695µs tracing=11.188µs
ts=2022-04-14T07:35:54.413Z caller=main.go:910 level=info msg="Server is ready to receive web requests."
```

The last line above confirms that the server has started successfully, which means our configuration is working correctly. However, starting the server manually is not ideal so we will create a new systemd unit file so that it can be managed through the Systemd daemon.

Kill the existing server with Ctrl-C, then create a prometheus.service file in the /etc/systemd/system directory as shown below:

```
sudo nano /etc/systemd/system/prometheus.service
```

Paste the following contents into the file and save it:

/etc/systemd/system/prometheus.service
```
[Unit]
Description=Prometheus Service
After=network.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

This service will perform a startup process similar to what we did manually before. It will also automatically restart the service if any errors occur and ensure the program is up and running after a system reboot.

Save and close the file, then reload the systemd configuration though the command below:

 
```
sudo systemctl daemon-reload
```

You can now start Prometheus by entering the following command:

 
```
sudo systemctl start prometheus
```

Go ahead and confirm that it started correctly:

```
systemctl status prometheus
```
You should observe the following output:

 Output
```
● prometheus.service - Prometheus Service
     Loaded: loaded (/etc/systemd/system/prometheus.service; disabled; vendor preset: enabled)
     Active: active (running) since Thu 2022-04-14 07:49:39 UTC; 8s ago
   Main PID: 1928667 (prometheus)
      Tasks: 6 (limit: 2275)
     Memory: 36.1M
     CGroup: /system.slice/prometheus.service
             └─1928667 /usr/local/bin/prometheus --config.file /etc/prometheus/prometheus.yml --storage.tsdb.path /var/lib/prometheus/ --web.console.templates=/etc/prometheus/consoles --web.>

Apr 14 07:49:39 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:39.939Z caller=head.go:536 level=info component=tsdb msg="On-disk memory mappable chunks replay completed" duration>
Apr 14 07:49:39 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:39.939Z caller=head.go:542 level=info component=tsdb msg="Replaying WAL, this may take a while"
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.280Z caller=head.go:613 level=info component=tsdb msg="WAL segment loaded" segment=0 maxSegment=1
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.281Z caller=head.go:613 level=info component=tsdb msg="WAL segment loaded" segment=1 maxSegment=1
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.282Z caller=head.go:619 level=info component=tsdb msg="WAL replay completed" checkpoint_replay_duration=61.49µs >
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.288Z caller=main.go:958 level=info fs_type=EXT4_SUPER_MAGIC
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.289Z caller=main.go:961 level=info msg="TSDB started"
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.289Z caller=main.go:1142 level=info msg="Loading configuration file" filename=/etc/prometheus/prometheus.yml
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.306Z caller=main.go:1179 level=info msg="Completed loading of configuration file" filename=/etc/prometheus/prome>
Apr 14 07:49:40 ubuntu-2gb-nbg1-1 prometheus[1928667]: ts=2022-04-14T07:49:40.307Z caller=main.go:910 level=info msg="Server is ready to receive web requests."
```


Finally, allow the autorun of the service: 

```
sudo systemctl enable prometheus
```

 Output
```
Created symlink /etc/systemd/system/multi-user.target.wants/prometheus.service → /etc/systemd/system/prometheus.service.
```

At this point, you have Prometheus running on the first system. In the next section, we will set up a Node Exporter instance on the second one.




### Step 6 — Downloading and installing Node Exporter
Just like Prometheus, Node Exporter cannot be installed from the default Ubuntu 20.04 repositories, so you need to get an archive of the current version on the Prometheus website . Copy the link to the appropriate Linux package, then log in to your second Linux server and download the archive:

Node Exporter Download Link
 
```
wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
```

You can also get the checksum of the downloaded file and compare it with the SHA256 Checksum similar to the way we confirmed the Prometheus download in step 3.

```
sha256sum node_exporter-1.3.1.linux-amd64.tar.gz
``` 

Output
```
68f3802c2dd3980667e4ba65ea2e1fb03f4a4ba026cca375f15a0390ff850949 node_exporter-1.3.1.linux-amd64.tar.gz
```

Next, unpack the archive and navigate to the directory with the unpacked files:

```
tar zxvf node_exporter-1.3.1.linux-amd64.tar.gz
```

Output
```
node_exporter-1.3.1.linux-amd64/
node_exporter-1.3.1.linux-amd64/LICENSE
node_exporter-1.3.1.linux-amd64/NOTICE
node_exporter-1.3.1.linux-amd64/node_exporter
``` 


```
cd node_exporter-1.3.1.linux-amd64
```
Afterward, copy the node_exporter binary to the /usr/local/bin directory, and set the user and group ownership to nodeuser:

 
```
sudo cp node_exporter /usr/local/bin/
sudo chown -R nodeuser:nodeuser /usr/local/bin/node_exporter
``` 

The Node Exporter CLI is now ready to go, but we need to secure it further so that only authorized clients can connect to it.




### Step 7 — Securing and running the Node Exporter instance
Before starting the Node Exporter, let's set up a secure connection. To do this, we will utilize an SSL connections, as well as basic authorization. First, generate a password hash through the htpasswd command:

```
htpasswd -nBC 10 "" | tr -d ':\n'
```
This command will prompt you to enter and confirm the password. Enter something you can remember, or use a password manager to generate one. For the purpose of this tutorial, we used the password prom. As a result, you will get a string like this:

 Output
```
$2y$10$fhJONsmSqbvjCB9sCjAe8.XyEDP39P36tCVoX8qx3sX3jcyCOhnHi
```

Copy the entire line above and write it down somewhere. It will be needed later in this step.

Now let's configure the Node Exporter by creating the following config file:

```
sudo nano /etc/node_exporter/web.yml
```

Paste the following contents into it:

```
tls_server_config:
  cert_file: /etc/node_exporter/node_exporter.crt
  key_file: /etc/node_exporter/node_exporter.key
basic_auth_users:
  prom: $2y$10$fhJONsmSqbvjCB9sCjAe8.XyEDP39P36tCVoX8qx3sX3jcyCOhnHi
```

We specified the path to the previously created certificate and the username and the hashed password required to connect. For demonstration purposes, we also used prom as the username.

Afterward, update the rights to the files:

```
sudo chown -R nodeuser:nodeuser /etc/node_exporter
```

To manage the Node Exporter instance with Systemd, we will create a new unit file like this:

```
sudo nano /etc/systemd/system/node_exporter.service
```
Paste the following contents into the file and save it.

/etc/systemd/system/node_exporter.service
```
[Unit]
Description=Node Exporter Service
After=network.target

[Service]
User=nodeuser
Group=nodeuser
Type=simple
ExecStart=/usr/local/bin/node_exporter --web.config=/etc/node_exporter/web.yml
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Next, reload the systemd configuration:

 
```
sudo systemctl daemon-reload

```

You can now start the Node Exporter service with systemctl:

 
```
sudo systemctl start node_exporter
```

Check that it started correctly:
 
```
systemctl status node_exporter
``` 

Output
```
● node_exporter.service - Node Exporter Service
     Loaded: loaded (/etc/systemd/system/node_exporter.service; disabled; vendor preset: disabled)
     Active: active (running) since Thu 2022-04-14 09:05:34 UTC; 6s ago
   Main PID: 302230 (node_exporter)
      Tasks: 4 (limit: 2254)
     Memory: 2.2M
        CPU: 15ms
     CGroup: /system.slice/node_exporter.service
             └─302230 /usr/local/bin/node_exporter --web.config=/etc/node_exporter/web.yml

Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=thermal_zone
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=time
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=timex
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=udp_queues
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=uname
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=vmstat
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.082Z caller=node_exporter.go:115 level=info collector=xfs
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.083Z caller=node_exporter.go:115 level=info collector=zfs
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.083Z caller=node_exporter.go:199 level=info msg="Listening on" address=:9100
Apr 14 09:05:34 fedora-2gb-hel1-1 node_exporter[302230]: ts=2022-04-14T09:05:34.085Z caller=tls_config.go:228 level=info msg="TLS is enabled." http2=true
```

Finally, enable Node Exporter to start on boot:

```
sudo systemctl enable node_exporter
```

Output

```
Created symlink /etc/systemd/system/multi-user.target.wants/node_exporter.service → /etc/systemd/system/node_exporter.service.
```

Before adding the newly created exporter to our Prometheus config, let's check that everything is working correctly. Return to the first server and use curl to make a request to the second server: 

```
curl https://<your_second_server_ip>:9100/metrics -k -u 'prom:prom'
```

where <your_second_server_ip> is the IP address of your second server.

You will observe some output that looks like like this:

 Output
```# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 8
. . .
```

If you don't get a similar output, check that the firewall is not blocking connections and also, check that the supplied username and password are correct.




### Step 8 — Configuring Prometheus to Scrape Node Exporter
In this step, we will update the Prometheus settings to be able to collect metrics from the Node Exporter. Staying on the first server, open the Prometheus config file:

 
```
sudo nano /etc/prometheus/prometheus.yml
```

Add the following lines to the end of the file under scrape_configs:

/etc/prometheus/prometheus.yml

```. . .
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]

  - job_name: 'node_exporter_client'
    scrape_interval: 5s
    scheme: https
    tls_config:
      ca_file: /etc/prometheus/node_exporter.crt
    basic_auth:
      username: prom
      password: prom
    static_configs:
      - targets: ['<your_second_server_ip>:9100']
```

The <your_second_server_ip> placeholder should be replaced with the IP address of your second server. Save the file and restart the Prometheus server:

```
sudo systemctl restart prometheus
```

Prometheus is now running and collecting information not only about itself but also from the specified remote server.




### Step 9 — Securing and working with the Prometheus web interface
Prometheus doesn't have any built-in authentication or authorization capabilities, so to prevent unauthorized access to data in Prometheus, we will use NGINX as a proxy. Ensure that it is installed on your first server:

 
```
sudo apt install nginx
```

If you are using a firewall you should also allow connections from NGINX:

 
```
sudo ufw allow 'Nginx HTTP'
```

Now, using the htpasswd utility, create a password file. In this case, we will use the username promadmin to access the web interface.

 
```
sudo htpasswd -c /etc/nginx/.htpasswd promadmin
```

Next, you need to set up NGINX. For our tutorial, we will use the default configuration file. Open it for editing:

 
```
sudo nano /etc/nginx/sites-enabled/default
```
Find the location / block under the server block. It should look like this:

/etc/nginx/sites-enabled/default
```
. . .
        location / {
                # First attempt to serve request as file, then
                # as directory, then fall back to displaying a 404.
                try_files $uri $uri/ =404;
        }
. . .
```

Replace this block with the following content:

/etc/nginx/sites-enabled/default
```
. . .
        location / {
          proxy_pass http://localhost:9090/;
          auth_basic "Prometheus";
          auth_basic_user_file "/etc/nginx/.htpasswd";
        }
. . .
```

Save and close the file, then check your NGINX configuration:

```
sudo nginx -t
```
The output should show that the syntax is ok and the test is successful:

 Output
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```
Afterward, reload the NGINX service:

```
sudo systemctl reload nginx
```

Open your web browser and navigate to http://<your_first_server_ip>. In the HTTP authentication dialogue box, enter the username and password you chose previously in this step. After a successful login, you will see the Prometheus UI: