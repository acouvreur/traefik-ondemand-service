# treafik-on-demand

## Description

This is a service that can scale up or down a docker swarm service on demand.
It basically starts a service when it's needed and then shut it down when it's no longer needed.

## Usage

In order to sue the service you should request the server according 
```
GET service_url/?name=<service_name>&timeout=<timeout>
```

`service_name`: The name of the service you want to call (and start if necessary)
`timeout`: The duration after which the service should be shut down if idle (in second)

Response:

`started`: The service is already started
`starting`: The service is starting


## Run 

To simply run the server you can use `go run main.go`.

## Deploy

To deploy this service in a container :

```
$ docker run -v /var/run/docker.sock:/var/run/docker.sock acouvreur/traefik-ondemand-service:latest
```
