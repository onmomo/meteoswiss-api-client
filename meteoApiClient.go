package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
)

func main() {
	postalCodePtr := flag.String("plz", "8000", "a CH postal code")
	hostPtr := flag.String("host", "10.0.1.2:9999", "the host that will receive the data")
	protocolPtr := flag.String("protocol", "udp", "the protocol for the host connection (tcp, udp and IP networks)")
	flag.Parse()
	initLogger()
	read(*postalCodePtr, *hostPtr, *protocolPtr)
}

func initLogger() {
	logger := log.Logger()
	logger.SetAppender(appenders.RollingFile("meteoswiss-api.log", true))
	appender := logger.Appender()
	appender.SetLayout(layout.Pattern("%d %p - %m%n"))
}

func read(postalCode string, host string, protocol string) {
	log.Info("Getting data for '%s' ...", postalCode)

	response, err := http.Get("http://www.meteoschweiz.admin.ch/product/output/danger/version__20180609_2026/de/dangers.json")
	if err != nil {
		log.Error("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Info(string(data))
	}

	log.Info("Closing Meteo Api Client, bye bye.")

	os.Exit(0)
}
