package client

import (
	"fmt"
	"io"
	"net/http"
)

func (client *SmartSchoolClient) CheckIfAuthenticated() error {
	if client.Pid == "" || client.PhpSessId == "" {
		return &AuthException{ApiException{"PID or PHPSESSID are not set"}}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/", client.domain), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s; pid=%s", client.PhpSessId, client.Pid))
	client.apiLogger.Info("Checking authentication status")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.apiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusFound {
		return &AuthException{ApiException{"Not authenticated, invalid cookies (PID or PHPSESSID)"}}
	} else if resp.StatusCode != http.StatusOK {
		return &ApiException{"Could not check if authenticated"}
	}

	return nil
}
