package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func NewSmartSchoolClient(domain string) *SmartSchoolClient {
	client := &SmartSchoolClient{
		domain: domain,

		WriteApiLogs:    os.Getenv("WRITE_API_LOGS") == "true",
		ApiLogger:       logrus.New(),
		WebsocketLogger: logrus.New(),
		AuthLogger:      logrus.New(),
	}

	client.ApiLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	client.WebsocketLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	client.AuthLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	if client.WriteApiLogs {
		err := os.Mkdir("./requests", 0755)
		if err != nil {
			client.ApiLogger.Error(fmt.Sprintf("Could not create ./requests directory: %s", err))
		}
	}

	return client
}
