package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gosmartschool/client/microsoftlogin"
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
		if _, err := os.Stat("./requests"); os.IsNotExist(err) {
			err := os.Mkdir("./requests", 0755)
			if err != nil {
				client.ApiLogger.Error(fmt.Sprintf("Could not create ./requests directory: %s", err))
			}
		}
	}

	return client
}

func (client *SmartSchoolClient) MicrosoftLogin(domain string, microsoftEmail string, microsoftPassword string, twoFactorSecurityQuestions microsoftlogin.TwoFactorSecurityQuestions) (bool, error) {
	phpSessId, err := microsoftlogin.MicrosoftLogin(domain, microsoftEmail, microsoftPassword, twoFactorSecurityQuestions)
	if err != nil {
		return false, err
	}

	client.PhpSessId = phpSessId
	err = client.CheckIfAuthenticated()
	if err != nil {
		return false, err
	}
	return true, nil
}
