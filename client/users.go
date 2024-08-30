package client

import (
	"encoding/xml"
	"fmt"
	"gosmartschool/structs"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FindUsersByName searches for users by name and returns a slice of structs.User.
func (client *SmartSchoolClient) FindUsersByName(name string) ([]structs.User, error) {
	client.apiLogger.Info("Requesting user from API")

	data := url.Values{}
	request := structs.SearchUsersRequest{
		Val:          name,
		Type:         0,
		ParentNodeID: "insertSearchFieldContainer_0_0",
		XML:          "<results></results>",
	}

	data.Set("val", request.Val)
	data.Set("type", fmt.Sprint(request.Type))
	data.Set("parentNodeId", request.ParentNodeID)
	data.Set("xml", request.XML)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/?module=Messages&file=searchUsers", client.domain), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("pid=%s; PHPSESSID=%s", client.Pid, client.PhpSessId))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.apiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		client.apiLogger.Info("User received")
		var response structs.SearchUsersResponse
		xmlData, _ := io.ReadAll(resp.Body)
		if err := xml.Unmarshal(xmlData, &response); err != nil {
			return nil, err
		}

		return response.Users, nil
	}

	client.apiLogger.Error("Could not get user")
	return nil, &ApiException{"Could not get user"}
}
