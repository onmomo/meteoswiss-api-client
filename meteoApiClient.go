package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
)

type Dangers struct {
	Config Config
	Days   map[string]Day
}

type Config struct {
	Name      string
	Language  string
	Version   string
	Timestamp Timestamp
}

type Day struct {
	Hazards   Hazards
	Lakes     []Hazard
	Airfields []Hazard
}

type Hazards struct {
	Wind          *Hazard
	Thunderstorm  *Hazard
	Snow          *Hazard
	Rain          *Hazard
	SlipperyRoads *Hazard
	HeatWave      *Hazard
	Frost         *Hazard
}

type Hazard struct {
	Description string
	Onset       Timestamp
	Expires     Timestamp
	Areas       []int
	IsOutlook   bool
	Warnlevel   int
	Name        string
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}

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

	// TODO use current timestamp
	response, err := http.Get("http://www.meteoschweiz.admin.ch/product/output/danger/version__20180609_2026/de/dangers.json")
	if err != nil {
		log.Error("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Info(string(data))
		var dangers Dangers
		json.Unmarshal([]byte(data), &dangers)

		log.Info(fmt.Sprintf("Warnlevel is: %d", dangers.Days["20180609_24h"].Lakes[0].Warnlevel))
	}

	log.Info("Closing Meteo Api Client, bye bye.")

	os.Exit(0)
}
