Robin
============
<img align="right" src="https://i.imgur.com/r2CWdQf.png" height="20%" width="20%">

Robin is a simple SSL termination server written in Go that allows you to serve your stuff off HTTPs while keeping 
downstream services on HTTP.

Robin is:
- A simple reverse proxy server
- A hot-headed vigilante with deep-rooted fears of a clown and crowbars.

Goals:
- [x] Provide SSL termination for multiple downstream services through a single endpoint
- [x] Provide an easy to use management API to control downstream services
- [ ] Make it work seamlessly in auto-scaling environments
- [ ] Provide a way to auto-assign SSL certificates from LetsEncrypt
- [ ] Somehow make it viable to use in production
- [x] Be free and open-source... Always.
- [x] Be a community-driven project.

Non-goals:
- Being a load balancer
- Being a WAF
- Being a certificate management service

Open TODOs:
- Don't half-ass the API
- Do better logging and error-handling
- Implement a pretty front-end sometime in the future
- Add tests for all possible packages
- Support HTTP to HTTPs redirection

## Proof of Concept

```shell
$ go test ./... -v
```
## How?

Edit `robinconfig.yaml` file to your liking and then `docker-compose build && docker-compose up`.
When running with compose, the traffic port and management port listens on 443 (HTTPS) and 8089 (HTTP) respectively.
You can map your DNS entries to the public address of the server running Robin. When an HTTPs resources is accessed,
Robin chooses an appropriate downstream server based on the hostname and routes to it.

An easy to use REST API is exposed under the management interface with these functions:

#### `GET /api/vhosts/` 
Lists configured hosts

Response:
```json
[
    {
        "id": 1,
        "created_at": "2020-06-10T18:23:39Z",
        "updated_at": "2020-06-10T18:23:39Z",
        "fqdn": "https://archlinux.localdomain",
        "origin": "http://localhost:8081",
        "certificate": {
            "id": 0,
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "rsa_key": null,
            "certificate": null,
            "ca_chain": null
        }
    }
]
```

#### `GET /api/vhosts/{id}`
Gets a single configured host

Response:
```json
{
    "id": 1,
    "created_at": "2020-06-10T18:23:39Z",
    "updated_at": "2020-06-10T18:23:39Z",
    "fqdn": "https://archlinux.localdomain",
    "origin": "http://localhost:8081",
    "certificate": {
        "id": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z",
        "rsa_key": "-----BEGIN PRIVATE KEY----- ......",
        "certificate": "-----BEGIN CERTIFICATE----- ......",
        "ca_chain": null
    }
}
```

#### `POST /api/vhosts/`
Creates a new host entry

Request:
```json
{
  "fqdn": "https://archlinux.localdomain",
  "origin": "http://someserver.com:8081",
  "cert": "-----BEGIN CERTIFICATE----- ......",
  "rsa": "-----BEGIN PRIVATE KEY----- ......"
}
```

Response: *same as get single host*

#### `DELETE /api/vhosts/{id}`
Deletes a host entry

> The management API uses basic authentication from the credentials configured in `robinconf.yaml`

### Command Line Args

- `--config`: Specify the configuration yaml file.
- `--logging-level`: Specify the logging level. Must be one of: CRITICAL, ERROR, WARNING, NOTICE, DEBUG, INFO
