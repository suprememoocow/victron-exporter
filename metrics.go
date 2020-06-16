package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "victron"

var (
	connectionStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mqtt_connection_state",
		Help:      "0=Disconnected; 1=Connected",
	}, []string{"client_id"})

	connectionStatusSinceTimeSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mqtt_connection_state_since_time_seconds",
		Help:      "Time since last change to mqtt_connection_state",
	}, []string{"client_id"})

	subscriptionsUpdatesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "mqtt_subscription_updates_total",
		Help:      "MQTT subscriptions updated received",
	})

	subscriptionsUpdatesIgnoredTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "mqtt_subscription_updates_ignored_total",
		Help:      "MQTT subscription updates ignored",
	})
)

func init() {
	prometheus.MustRegister(connectionStatus)
	prometheus.MustRegister(connectionStatusSinceTimeSeconds)
	prometheus.MustRegister(subscriptionsUpdatesTotal)
	prometheus.MustRegister(subscriptionsUpdatesIgnoredTotal)
}
