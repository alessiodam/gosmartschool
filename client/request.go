package client

import (
	"fmt"
	"net/http"
	"strings"
)

func (client *SmartSchoolClient) sendXmlRequest(method string, url string, body string, extraHeaders map[string]string) (*http.Response, error) {
	client.apiLogger.Info("Sending request to API")

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s%s", client.domain, url), strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("pid=%s; PHPSESSID=%s", client.Pid, client.PhpSessId))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	client.apiLogger.Info("Response Status Code from API: %s", resp.StatusCode)

	return resp, nil
}
