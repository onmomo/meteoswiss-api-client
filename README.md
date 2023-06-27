# meteoswiss-api-client
Consumes weather forecasts and hazards via Meteoswiss api for a specific postal code.
In case of an acute hazard alert for the specified postal code, a configured target host will be notified.

## How to use
To check for thunderstorm reports for postalcode `9500` and write the payload `hazard:1` to the `targetHost` every 10 minutes run the following command

```
docker run --rm onmomo/meteoswiss-api-client:main --plz 9500 --targetHost 10.0.1.99:8080
```

### Configuration
To get a list off all available parameters, run
```
docker run --rm onmomo/meteoswiss-api-client:main --help

Usage of ./app:
  -cron string
    	cron expression to configure the hazard poll interval. Defaults to every 10 minutes (default "0 0/10 * * * *")
  -plz string
    	a CH postal code (default "9500")
  -protocol string
    	the protocol for the target host connection (tcp, udp and IP networks) (default "udp")
  -targetHost string
    	the host that will receive the data (default "10.0.1.2:9990")
```


## How to develop

build and run using docker
```
docker build -t onmomo/meteoswiss-api-client:local .
docker run --rm onmomo/meteoswiss-api-client:local
```
or without docker
```
go build -v
./meteoswiss-api-client --plz 9500 --targetHost 10.0.1.2:8888
```
