package client

import (
	"io"
	"net/http"
)

func (client *SmartSchoolClient) CheckIfAuthenticated() error {
	if client.Pid == "" || client.PhpSessId == "" {
		return &AuthException{ApiException{"PID or PHPSESSID are not set"}}
	}

	client.apiLogger.Info("Checking authentication status")
	resp, err := client.sendXmlRequest("GET", "/", "", nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.apiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusUnauthorized {
		return &AuthException{ApiException{"Not authenticated, invalid cookies (PID or PHPSESSID)"}}
	} else if resp.StatusCode != http.StatusOK {
		return &ApiException{"Could not check if authenticated"}
	}

	return nil
}
