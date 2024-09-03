package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (client *SmartSchoolClient) sendRequest(method string, url string, body string, extraHeaders *http.Header) (*http.Response, string, error) {
	client.ApiLogger.Info(fmt.Sprintf("Sending request to API: %s %s", method, url))

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s%s", client.domain, url), strings.NewReader(body))
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Cookie", fmt.Sprintf("pid=%s; PHPSESSID=%s", client.Pid, client.PhpSessId))
	req.Header.Set("Host", client.domain)
	req.Header.Set("Origin", fmt.Sprintf("https://%s", client.domain))

	if extraHeaders != nil {
		for k, v := range *extraHeaders {
			req.Header[k] = v
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("./requests/%s.txt", timestamp)

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			client.ApiLogger.Error(err)
		}
	}(file)

	_, _ = fmt.Fprintf(file, "Request:\n")
	_, _ = fmt.Fprintf(file, "Method: %s\n", method)
	_, _ = fmt.Fprintf(file, "URL: %s\n", fmt.Sprintf("https://%s%s", client.domain, url))
	_, _ = fmt.Fprintf(file, "Headers:\n")
	for k, v := range req.Header {
		_, _ = fmt.Fprintf(file, "%s: %s\n", k, strings.Join(v, ", "))
	}
	_, _ = fmt.Fprintf(file, "Body:\n%s\n\n", body)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	_, _ = fmt.Fprintf(file, "Response Status Code: %d\n", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	_, _ = fmt.Fprintf(file, "Response Body:\n%s\n", string(respBody))

	client.ApiLogger.Info(fmt.Sprintf("Response Status Code from API: %d", resp.StatusCode))

	return resp, string(respBody), nil
}

func (client *SmartSchoolClient) sendXmlRequest(method string, url string, body string, extraHeaders *http.Header) (*http.Response, string, error) {
	if extraHeaders == nil {
		extraHeaders = &http.Header{}
	}
	extraHeaders.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	extraHeaders.Set("X-Requested-With", "XMLHttpRequest")
	return client.sendRequest(method, url, body, extraHeaders)
}
