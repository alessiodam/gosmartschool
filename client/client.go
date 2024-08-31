package client

import (
	"github.com/sirupsen/logrus"
)

func NewSmartSchoolClient(domain string) *SmartSchoolClient {
	client := &SmartSchoolClient{
		domain:          domain,
		apiLogger:       logrus.New(),
		websocketLogger: logrus.New(),
		authLogger:      logrus.New(),
	}

	client.apiLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	client.websocketLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	client.authLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	return client
}
