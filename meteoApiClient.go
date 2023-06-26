package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
	"github.com/procyon-projects/chrono"
)

type Weather struct {
	Warnings []Warnings
}

type Warnings struct {
	HtmlText  string
	Ordering  string
	Outlook   bool
	Text      string
	ValidFrom Timestamp
	ValidTo   Timestamp
	WarnLevel int
	WarnType  WarnType
}

type WarnType int

// id 1 = thunderstorm
// id 2 = rain
// id 10 = forecast fire
// id 11 = flood
// id xx = wind
// id xx = slippery-roads
// id xx = frost
// id xx = heat-wave
const (
	Thunderstorm WarnType = 1
	Rain         WarnType = 2
	Flood        WarnType = 11
	ForestFire   WarnType = 10
)

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

var wg sync.WaitGroup

func main() {
	postalCodePtr := flag.String("plz", "9500", "a CH postal code")
	hostPtr := flag.String("targetHost", "10.0.1.2:9990", "the host that will receive the data")
	protocolPtr := flag.String("protocol", "udp", "the protocol for the target host connection (tcp, udp and IP networks)")
	cronExpressionPtr := flag.String("cron", "0 0/10 * * * *", "cron expression to configure the hazard poll interval. Defaults to every 10 minutes")
	flag.Parse()
	initLogger()

	wg.Add(1)
	taskScheduler := chrono.NewDefaultTaskScheduler()
	_, err := taskScheduler.ScheduleWithCron(func(ctx context.Context) {
		_, err := read(*postalCodePtr, *hostPtr, *protocolPtr)
		if err != nil {
			log.Warn("Meteoswiss Api Client scheduled run failed", err)
		}
		log.Info("Meteoswiss Api Client task ran successfully, bye bye.")
	}, *cronExpressionPtr)

	if err == nil {
		log.Info("Meteoswiss Api Client task has been scheduled ..")
	} else {
		log.Warn("Meteoswiss Api Client scheduled task initialization failed", err)
	}
	wg.Wait()
}

func initLogger() {
	logger := log.Logger()
	logger.SetAppender(appenders.Console())
	appender := logger.Appender()
	appender.SetLayout(layout.Pattern("%d %p - %m%n"))
}

func read(postalCode string, host string, protocol string) (bool, error) {
	weatherByPostalCodeAPI, err := resolveWeatherByPostalCodeAPI(postalCode)
	if err != nil {
		return false, fmt.Errorf("could not resolve weatherByPostalCodeApi %s", err)
	}
	log.Debug("Requesting weather from %s ...", weatherByPostalCodeAPI)
	weatherResponse, err := http.Get(weatherByPostalCodeAPI)
	if err != nil {
		return false, fmt.Errorf("the HTTP request failed with error %s", err)
	} else if weatherResponse.StatusCode != 200 {
		return false, fmt.Errorf("unable to retrieve the weather for postalcode '%s'", postalCode)
	} else {
		data, _ := io.ReadAll(weatherResponse.Body)
		log.Debug(fmt.Sprintf("Weather JSON response: %s", data))
		var weather Weather
		json.Unmarshal([]byte(data), &weather)

		if len(weather.Warnings) > 0 {
			log.Info("Opening %s connection to %s ...", protocol, host)
			conn, err := net.Dial(protocol, host)
			if err != nil {
				log.Fatal(err)
				return false, fmt.Errorf("couldn't open %s connection to %s", protocol, host)
			}

			log.Info(fmt.Sprintf("Following warnings were reported for postal code: %s", postalCode))
			for _, warning := range weather.Warnings {
				log.Info(fmt.Sprintf("Warnlevel: %d, Warntype: %d, Outlook: %t: %s", warning.WarnLevel, warning.WarnType, warning.Outlook, warning.Text))
				if !warning.Outlook && warning.WarnLevel >= 3 && warning.WarnType == Thunderstorm {
					log.Info(fmt.Sprintf("Send type %d alert, level %d.", warning.WarnType, warning.WarnLevel))
					conn.Write([]byte("hazard:1"))
				}
			}

			conn.Close()

			return true, nil
		} else {
			log.Info("No warnings found for postal code %s, all good.", postalCode)
			return true, nil
		}
	}
}

func resolveWeatherByPostalCodeAPI(postalCode string) (string, error) {
	weatherByPostalCode, err := url.Parse("https://app-prod-ws.meteoswiss-app.ch/v1/plzDetail")
	if err != nil {
		return "", err
	}
	query := weatherByPostalCode.Query()
	query.Set("plz", postalCode+"00")
	weatherByPostalCode.RawQuery = query.Encode()

	return weatherByPostalCode.String(), nil
}
