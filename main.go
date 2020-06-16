package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	)
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*listenAddress, nil)

	mqttOpts := mqttConnnectionConfig{*host, *port, *secure, *username, *password}
	go func() {
		err := listen(*clientPrefix+"_sub", mqttOpts, "#")
		if err != nil {
			log.Fatal(err)
		}
	}()

	client, err := connect(*clientPrefix+"_pub", mqttOpts)
	if err != nil {
		log.Fatal(err)
	}

	timer := time.NewTicker(*pollInterval)
	for range timer.C {
		if !client.IsConnectionOpen() || systemSerialID == "" {
			continue
		}

		token := client.Publish(fmt.Sprintf("R/%s/system/0/Serial", systemSerialID), 1, false, "")
		for !token.WaitTimeout(5 * time.Second) {
			if err := token.Error(); err != nil {
				log.Printf("mqtt publish failed: %v", err)
			}
		}
	}
}
