package main

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			log.WithFields(log.Fields{
				"env_var":   key,
				"env_value": value}).
				WithError(err).Fatal("Unable to parse ENV VAR as an INT")
		}
		return i
	}
	log.WithFields(log.Fields{
		"env_var":        key,
		"fallback_value": fallback}).
		Debug("Unable to find ENV VAR, falling back to default value")

	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			log.WithFields(log.Fields{
				"env_var":   key,
				"env_value": value}).
				WithError(err).Fatal("Unable to parse ENV VAR as a BOOL")
		}
		return b
	}
	log.WithFields(log.Fields{
		"env_var":        key,
		"fallback_value": fallback}).
		Debug("Unable to find ENV VAR, falling back to default value")

	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err != nil {
			log.WithFields(log.Fields{
				"env_var":   key,
				"env_value": value}).
				WithError(err).Fatal("Unable to parse ENV VAR as a DURATION")
		}
		return d
	}
	log.WithFields(log.Fields{
		"env_var":        key,
		"fallback_value": fallback}).
		Debug("Unable to find ENV VAR, falling back to default value")

	return fallback
}
