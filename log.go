package main

import (
	log "github.com/sirupsen/logrus"
)

func setLogLevel(logLevel int) {
	switch logLevel {
	case 0:
		log.SetLevel(log.DebugLevel)
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.WarnLevel)
	case 3:
		log.SetLevel(log.ErrorLevel)
	case 4:
		log.SetLevel(log.FatalLevel)
	}
}
