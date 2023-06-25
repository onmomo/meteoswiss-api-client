# meteoswiss-api-client
Consumes weather forecasts and hazards via Meteoswiss api for a specific postal code.
In case of an acute hazard alert for the specified postal code, a configured target host will be notified.

## How to develop

build the project
```
go build -v
```
and run it
```
./meteoswiss-api-client --plz 9500 --targetHost 10.0.1.2:8888
```
