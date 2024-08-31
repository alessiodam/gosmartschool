package client

import (
	"encoding/xml"
	"fmt"
	"gosmartschool/structs"
	"io"
	"net/http"
	"net/url"
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

	resp, err := client.sendXmlRequest("POST", "/?module=Messages&file=searchUsers", data.Encode(), nil)
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
