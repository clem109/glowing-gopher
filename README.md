[![Go Report Card](https://goreportcard.com/badge/clem109/glowing-gopher)](https://goreportcard.com/report/clem109/glowing-gopher)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/clem109/glowing-gopher/blob/master/LICENSE)
# Glowing-Gopher

<img src="../master/gopher.jpeg?raw=true" width="300" height="300" />

This project is to check the health status of different microservices, simply clone this repo and edit the config.yml file to your needs.

The default port is ":3333".

Call the "/healthcheck" endpoint and it will automatically test all the endpoints provided and return their status codes.

### The current binary is with the existing config.yml so this needs to be rebuilt with your own config file.

## Build

Build the binary with the following commands:

```
go build main.go
go install
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=linux CGO_ENABLED=0 go build main.go
docker build --rm -t clem109/glowing-gopher .
docker run -it -p [YOUR PORT]:3333 clem109/glowing-gopher
```

In our dockerfile we use the compiled binary as it makes the image size much much smaller.
