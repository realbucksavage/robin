Robin
============
<img align="right" src="https://i.imgur.com/r2CWdQf.png" height="20%" width="20%">

Robin is a simple SSL termination server written in Go. it has the following features.

- Add/remove new hosts on runtime
- Generate and apply certificates from LetsEncrypt
- Upload custom certificates
- Helps Batman save Gotham.

## How?

Edit `robinconfig.yaml` file to your liking and then `docker-compose build && docker-compose up`.
When running with compose, the traffic port and management port listens on 443 (HTTPS) and 8089 (HTTP) respectively.

### Command Line Args

- `--config`: Specify the configuration yaml file.
- `--logging-level`: Specify the logging level. Must be one of: CRITICAL, ERROR, WARNING, NOTICE, DEBUG, INFO
