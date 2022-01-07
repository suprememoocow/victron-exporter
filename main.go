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

var (
	listenAddress = flag.String("web.listen-address",
		getEnv("LISTEN_ADDR", "127.0.0.1:9226"),
		"Address on which to expose metrics and web interface.")

	clientPrefix = flag.String("mqtt.client_prefix",
		getEnv("MQTT_CLIENT_PREFIX", "victron_exporter"),
		"Prefix for MQTT clientID")

	pollInterval = flag.Duration("victron.poll_interval",
		getDurationEnv("VICTRON_POLL_INTERVAL", 10*time.Second),
		"CGX MQTT poll interval")

	host = flag.String("mqtt.host",
		getEnv("MQTT_HOST", ""),
		"CGX IP address or hostname")

	port = flag.Int("mqtt.port",
		getIntEnv("MQTT_PORT", 8883),
		"CGX MQTT PORT")

	secure = flag.Bool("mqtt.secure",
		getBoolEnv("MQTT_SECURE", true),
		"CGX SSL-enabled communication")

	username = flag.String("mqtt.username",
		getEnv("MQTT_USERNAME", ""),
		"Victron MQTT Cloud Username")

	password = flag.String("mqtt.password",
		getEnv("MQTT_PASSWORD", ""),
		"Victron MQTT Cloud Password")

	logLevel = flag.Int("log.level",
		getIntEnv("LOG_LEVEL", 2),
		"Log level: 0=debug, 1=info, 2=warn, 3=error")
)

func main() {
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
