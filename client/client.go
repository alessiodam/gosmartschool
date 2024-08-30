package client

import (
	"github.com/sirupsen/logrus"
)

type ApiException struct {
	Message string
}

func (e *ApiException) Error() string {
	return e.Message
}

type AuthException struct {
	ApiException
}

type SmartSchoolClient struct {
	domain                  string
	platformId              string
	PhpSessId               string
	Pid                     string
	userID                  string
	receivedMessageCallback func(message map[string]interface{})
	userToken               string
	apiLogger               *logrus.Logger
	websocketLogger         *logrus.Logger
	authLogger              *logrus.Logger
}

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
