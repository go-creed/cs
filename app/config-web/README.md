# Config Service

This is the Config service

Generated with

```
micro new --namespace=go.micro.cs --type=web cs/app/config-web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.cs.web.config
- Type: web
- Alias: config

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
go run main.go --registry=etcd --registry_address=127.0.0.1:2379 --cc=127.0.0.1:2379
```

Build a docker image
```
make docker
```