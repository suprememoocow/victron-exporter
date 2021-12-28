package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var systemSerialID = ""

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9226",
			"Address on which to expose metrics and web interface.")
		host         = flag.String("mqtt.host", "", "")
		port         = flag.Int("mqtt.port", 8883, "")
		secure       = flag.Bool("mqtt.secure", true, "")
		clientPrefix = flag.String("mqtt.client_prefix", "victron_exporter", "Prefix for MQTT clientID")
		username     = flag.String("mqtt.username", "", "")
		password     = flag.String("mqtt.password", "", "")
		pollInterval = flag.Duration("victron.poll_interval", 10*time.Second, "MQTT poll interval")
		logLevel     = flag.Int("log.level", 2, "Log level: 0=debug, 1=info, 2=warn, 3=error")
	)
	flag.Parse()

	setLogLevel(*logLevel)

	log.WithField("address", *listenAddress).Info("victron_exporter listening")

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(*listenAddress, nil)
		if err != nil {
			log.WithField("address", *listenAddress).WithError(err).Fatal("failed to listen on address")
		}
	}()

	mqttOpts := mqttConnectionConfig{*host, *port, *secure, *username, *password}
	go func() {
		err := listen(*clientPrefix+"_sub", mqttOpts, "#")
		if err != nil {
			log.WithError(err).Fatal("failed to establish mqtt subscription connection")
		}
	}()

	client, err := connect(*clientPrefix+"_pub", mqttOpts)
	if err != nil {
		log.WithError(err).Fatal("failed to establish mqtt publish connection")
	}

	timer := time.NewTicker(*pollInterval)
	for range timer.C {
		if !client.IsConnectionOpen() {
			log.Debug("mqtt connection not yet established")

			continue
		}

		// Check whether we've heard back from victron mqtt yet...
		if systemSerialID == "" {
			log.Debug("awaiting system serial ID response from Victron mqtt bus")

			continue
		}

		token := client.Publish(fmt.Sprintf("R/%s/system/0/Serial", systemSerialID), 1, false, "")
		for !token.WaitTimeout(5 * time.Second) {
			if err := token.Error(); err != nil {
				log.WithError(err).Error("mqtt publish failed")
			}
		}
	}
}
