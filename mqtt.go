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
		InsecureSkipVerify: true,
		// // Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
	}
}

func connect(mqttOptions *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(mqttOptions)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}

func createClientOptions(clientID string, host string, port int, secure bool, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	if secure {
		opts.AddBroker(fmt.Sprintf("ssl://%s:%d", host, port))
		opts.SetTLSConfig(newTLSConfig())
	} else {
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	}

	if username != "" {
		opts.SetUsername(username)
	}

	if password != "" {
		opts.SetPassword(password)
	}

	opts.SetClientID(clientID)
	return opts
}

type victronValue struct {
	Value *float64 `json:"value"`
}

type victronStringValue struct {
	Value *string `json:"value"`
}

func listen(mqttOptions *mqtt.ClientOptions, topic string) error {
	client, err := connect(mqttOptions)
	if err != nil {
		return err
	}
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		topicInfoParts := topicParts[4:]

		componentType := topicParts[2]
		componentID := topicParts[3]

		topicString := strings.Join(topicInfoParts, "/")

		if topicString == "Serial" {
			var v victronStringValue
			json.Unmarshal(msg.Payload(), &v)
			log.Println("SSetting serial to " + *v.Value)
			systemSerialID = *v.Value
		}

		o, ok := suffixTopicMap[topicString]
		if ok {
			var v victronValue
			json.Unmarshal(msg.Payload(), &v)

			if v.Value == nil {
				o(componentType, componentID, math.NaN())
			} else {
				o(componentType, componentID, *v.Value)
			}
		}
	})

	return nil
}
