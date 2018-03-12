# Glowing-Gopher

[[https://github.com/clem109/repository/blob/master/gopher.jpeg|alt=glowing-gopher]]

This project is to check the health status of our services, it will notify us whether there are any issues with any of our services that we add and provides an API for testing whether the service itself is running.

It's written in Go because it's a small project and Go is cool.

A yaml file is all that is needed to get up and running, the yaml file will include information about the endpoint to hit, expected response, and the frequency of API requests (in a cron job format). We will also expose an API endpoint to trigger healthchecks on all our services whenever requested.
