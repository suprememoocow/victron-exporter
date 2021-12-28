package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func newTLSConfig() *tls.Config {
	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: roots,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true, //nolint:gosec
		// // Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
	}
}

type mqttConnnectionConfig struct {
	host     string
	port     int
	secure   bool
	username string
	password string
}

func connectWait(client mqtt.Client) error {
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}

	err := token.Error()
	if err != nil {
		return fmt.Errorf("failed to connect to mqtt: %w", err)
	}

	return nil
}

func connect(clientID string, config mqttConnnectionConfig) (mqtt.Client, error) {
	client := mqtt.NewClient(createClientOptions(clientID, config, nil))

	return client, connectWait(client)
}

func listen(clientID string, config mqttConnnectionConfig, topic string) error {
	onConnect := func(client mqtt.Client) {
		// We need to subscribe after each connection
		// since mqtt does not maintain subscriptions across reconnects
		client.Subscribe(topic, 0, mqttSubscriptionHandler)
	}

	client := mqtt.NewClient(createClientOptions(clientID, config, onConnect))

	return connectWait(client)
}

func newConnectionLostHandler(clientID string) mqtt.ConnectionLostHandler {
	return func(c mqtt.Client, e error) {
		log.Printf("mqtt connection lost. clientId=%v, error=%v", clientID, e)
		connectionStatus.WithLabelValues(clientID).Set(0)
		connectionStatusSinceTimeSeconds.WithLabelValues(clientID).Set(float64(time.Now().Unix()))
	}
}

func newConnectionHandler(clientID string, wrapped mqtt.OnConnectHandler) mqtt.OnConnectHandler {
	return func(c mqtt.Client) {
		log.Printf("mqtt connected. clientId=%v", clientID)
		connectionStatus.WithLabelValues(clientID).Set(1)
		connectionStatusSinceTimeSeconds.WithLabelValues(clientID).Set(float64(time.Now().Unix()))

		if wrapped != nil {
			wrapped(c)
		}
	}
}

func createClientOptions(clientID string, config mqttConnnectionConfig, onConnectionHandler mqtt.OnConnectHandler) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)
	opts.SetWriteTimeout(30 * time.Second)
	opts.SetOrderMatters(false)
	opts.SetConnectionLostHandler(newConnectionLostHandler(clientID))
	opts.SetOnConnectHandler(newConnectionHandler(clientID, onConnectionHandler))

	if config.secure {
		opts.AddBroker(fmt.Sprintf("ssl://%s:%d", config.host, config.port))
		opts.SetTLSConfig(newTLSConfig())
	} else {
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.host, config.port))
	}

	if config.username != "" {
		opts.SetUsername(config.username)
	}

	if config.password != "" {
		opts.SetPassword(config.password)
	}

	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	return opts
}

type victronValue struct {
	Value *float64 `json:"value"`
}

type victronStringValue struct {
	Value *string `json:"value"`
}

func mqttSubscriptionHandler(client mqtt.Client, msg mqtt.Message) {
	subscriptionsUpdatesTotal.Inc()

	topic := msg.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) < 5 {
		subscriptionsUpdatesIgnoredTotal.Inc()

		return
	}
	topicInfoParts := topicParts[4:]

	componentType := topicParts[2]
	componentID := topicParts[3]

	topicString := strings.Join(topicInfoParts, "/")

	if (topicString == "Serial") && (systemSerialID == "") {
		var v victronStringValue

		err := json.Unmarshal(msg.Payload(), &v)
		if err != nil {
			subscriptionsUpdatesIgnoredTotal.Inc()

			return
		}

		if v.Value != nil {
			systemSerialID = *v.Value
		}

		return
	}

	o, ok := suffixTopicMap[topicString]
	if !ok {
		subscriptionsUpdatesIgnoredTotal.Inc()

		return
	}

	var v victronValue

	err := json.Unmarshal(msg.Payload(), &v)
	if err != nil {
		subscriptionsUpdatesIgnoredTotal.Inc()

		return
	}

	if v.Value == nil {
		o(componentType, componentID, math.NaN())
	} else {
		o(componentType, componentID, *v.Value)
	}
}
