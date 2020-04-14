# Upload Service

This is the Upload service

Generated with

```
micro new cs/app/upload-web --namespace=go.micro.cs --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.cs.web.upload
- Type: web
- Alias: upload

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
./upload-web
```

Build a docker image
```
make docker
```