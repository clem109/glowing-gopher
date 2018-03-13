# Glowing-Gopher

<img src="../master/gopher.jpeg?raw=true" width="300" height="300" />

This project is to check the health status of different microservices, simply clone this repo and edit the config.yml file to your needs.

The default port is ":3333", but this can be changed using the PORT envvar.

Call the "/healthcheck" endpoint and it will automatically test all the endpoints provided and return their status codes.
