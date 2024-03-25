# Teonet fortune web-server microservice

This is simple [Teonet](https://github.com/teonet-go/teonet) web-server microservice application which get fortune message from [Teonet Fortune](https://github.com/teonet-go/teofortune) microservice and show it in the site web page.

[![GoDoc](https://godoc.org/github.com/teonet-go/teofortune-web?status.svg)](https://godoc.org/github.com/teonet-go/teofortune-web/)
[![Go Report Card](https://goreportcard.com/badge/github.com/teonet-go/teofortune-web)](https://goreportcard.com/report/github.com/teonet-go/teofortune-web)

To create your own web-site you need to have host with IP address, domain and possibility to create record in this domain. We have created the [fortune.teonet.dev](https://fortune.teonet.dev) web-site and run it in Teonet Cloud.

## Run the Teonet fortune web-site microservice

There are various ways to start this Teonet microservice application:

_In code blow we use preinstalled in Teonet Cloud teofortune microservice address: '-fortune=8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF'.
Change this address to your application address. Or you can use this address, but than you will connect to `teofortune` microservice application running in Teonet Cloud. The address prints after you start Teonet application in string:_  
`Teonet address: 8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF`

### 1. From sources

```bash
git clone https://github.com/teonet-go/teofortune-web
cd teofortune-web
go run . -fortune=8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF -loglevel=debug
```

### 2. Install binary from github

```bash
go install github.com/teonet-go/teofortune-web .
teofortune -fortune=8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF -loglevel=debug
```

### 3. Using docker

```bash
docker run -d -it --network=host --restart=always --name teofortune-web -v \
$HOME/.config/teonet/teofortune-web:/root/.config/teonet/teofortune-web ghcr.io/teonet-go/\
teofortune-web:latest teofortune-web \
-fortune=8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF -loglevel=debug
```

## How to use

By default the teofortune-web site start at [localhost:8088](http://localhost:8088). If you run your web site in host with real IP you can add parameter `-domain=my.feature.example` to start real internet web server. Where the `my.feature.example` is your domain and you can create records in this domain.

There is preinstalled teofortune-web web-site with name [fortune.teonet.dev](https://fortune.teonet.dev)

## Licence

[BSD](LICENSE)
