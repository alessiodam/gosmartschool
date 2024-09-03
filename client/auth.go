package client

import (
	"io"
	"net/http"
)

func (client *SmartSchoolClient) CheckIfAuthenticated() error {
	if client.PhpSessId == "" {
		return &AuthException{ApiException{"PHPSESSID are not set"}}
	}

	client.ApiLogger.Info("Checking authentication status")
	resp, _, err := client.sendXmlRequest("GET", "/", "", nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.ApiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusUnauthorized {
		return &AuthException{ApiException{"Not authenticated, invalid cookies (PID or PHPSESSID)"}}
	} else if resp.StatusCode != http.StatusOK {
		return &ApiException{"Could not check if authenticated"}
	}

	return nil
}
