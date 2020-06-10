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
		host     = flag.String("mqtt.host", "", "")
		port     = flag.Int("mqtt.port", 8883, "")
		secure   = flag.Bool("mqtt.secure", true, "")
		username = flag.String("mqtt.username", "", "")
		password = flag.String("mqtt.password", "", "")
	)
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*listenAddress, nil)

	listenOpts := createClientOptions("victron_exporter_sub", *host, *port, *secure, *username, *password)
	go func() {
		err := listen(listenOpts, "#")
		if err != nil {
			log.Fatal(err)
		}
	}()

	publishClientOpts := createClientOptions("victron_exporter_pub", *host, *port, *secure, *username, *password)
	client, err := connect(publishClientOpts)
	if err != nil {
		log.Fatal(err)
	}

	timer := time.NewTicker(5 * time.Second)
	for range timer.C {
		if systemSerialID != "" {
			client.Publish(fmt.Sprint("R/%s/system/0/Serial", systemSerialID), 1, false, "")
		}
	}
}
