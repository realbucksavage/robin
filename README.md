Robin
============
<img align="right" src="https://i.imgur.com/r2CWdQf.png" height="20%" width="20%">

Robin is a simple SSL termination server written in Go that allows you to serve your stuff off HTTPs while keeping 
downstream services on HTTP.

Robin is:
- A simple reverse proxy server
- A hot-headed vigilante with deep-rooted fears of a clown and crowbars.

Robin is NOT:
- A load-balancer 
- A WAF of any kind

Goals:
- Provide SSL termination for multiple downstream services through a single endpoint
- Make it work seamlessly in auto-scaling environments
- Provide a way to auto-assign SSL certificates from LetsEncrypt.
- Implement a pretty front-end sometime in the future.
- Be free and open-source... Always.

## How?

Edit `robinconfig.yaml` file to your liking and then `docker-compose build && docker-compose up`.
When running with compose, the traffic port and management port listens on 443 (HTTPS) and 8089 (HTTP) respectively.

### Command Line Args

- `--config`: Specify the configuration yaml file.
- `--logging-level`: Specify the logging level. Must be one of: CRITICAL, ERROR, WARNING, NOTICE, DEBUG, INFO
