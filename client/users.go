package client

import (
	"encoding/xml"
	"gosmartschool/structs"
	"io"
	"net/http"
	"net/url"
)

// FindUsersByName searches for users by name and returns a slice of structs.User.
func (client *SmartSchoolClient) FindUsersByName(name string) ([]structs.User, error) {
	client.ApiLogger.Info("Requesting user from API")

	data := url.Values{}
	data.Set("val", name)
	data.Set("type", "0")
	data.Set("parentNodeId", "insertSearchFieldContainer_0_0")
	data.Set("xml", "<results></results>")
	data.Set("uniqueUsc", client.UniqueUsc)

	resp, _, err := client.sendXmlRequest("POST", "/?module=Messages&file=searchUsers", data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			client.ApiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		client.ApiLogger.Info("User received")
		var response struct {
			Users []structs.User `xml:"users>user"`
		}

		xmlData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Adding <results> wrapper to match the expected XML structure
		xmlData = []byte("<results>" + string(xmlData) + "</results>")

		if err := xml.Unmarshal(xmlData, &response); err != nil {
			return nil, err
		}

		return response.Users, nil
	}

	client.ApiLogger.Error("Could not get user")
	return nil, &ApiException{"Could not get user"}
}
