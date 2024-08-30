package client

import (
	"encoding/xml"
	"fmt"
	"gosmartschool/structs"
	"io"
	"net/http"
	"net/url"
)

// parseMessageListResponse parses the XML response for a list of messages and returns a slice of structs.Message.
func parseMessageListResponse(body io.Reader) ([]structs.Message, error) {
	var response structs.MessageListResponse
	if err := xml.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}

// ListMessages retrieves a list of messages from the API and returns a slice of structs.Message.
func (client *SmartSchoolClient) ListMessages() ([]structs.Message, error) {
	client.apiLogger.Info("Requesting messages from API")

	data := url.Values{}
	data.Set("command", `
		<request>
			<command>
				<subsystem>postboxes</subsystem>
				<action>message list</action>
				<params>
					<param name="boxType"><![CDATA[inbox]]></param>
					<param name="boxID"><![CDATA[0]]></param>
					<param name="sortField"><![CDATA[date]]></param>
					<param name="sortKey"><![CDATA[desc]]></param>
					<param name="poll"><![CDATA[false]]></param>
					<param name="poll_ids"><![CDATA[]]></param>
					<param name="layout"><![CDATA[new]]></param>
				</params>
			</command>
		</request>
	`)

	resp, err := client.sendXmlRequest("POST", "/?module=Messages&file=dispatcher", data.Encode(), nil)
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
		client.apiLogger.Info("Messages received")
		return parseMessageListResponse(resp.Body)
	}

	client.apiLogger.Error("Could not get messages")
	return nil, &ApiException{"Could not get messages"}
}

// GetMessageByID retrieves a single message by its ID and returns a structs.Message.
func (client *SmartSchoolClient) GetMessageByID(messageID string) (structs.Message, error) {
	client.apiLogger.Info("Requesting message from API")

	data := url.Values{}
	data.Set("command", fmt.Sprintf(`
		<request>
			<command>
				<subsystem>postboxes</subsystem>
				<action>show message</action>
				<params>
					<param name="msgID"><![CDATA[%s]]></param>
					<param name="boxType"><![CDATA[inbox]]></param>
					<param name="limitList"><![CDATA[true]]></param>
				</params>
			</command>
		</request>
	`, messageID))

	resp, err := client.sendXmlRequest("POST", "/?module=Messages&file=dispatcher", data.Encode(), nil)
	if err != nil {
		return structs.Message{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.apiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		client.apiLogger.Info("Message received")
		return parseSingleMessageResponse(resp.Body)
	}

	client.apiLogger.Error("Could not get message")
	return structs.Message{}, &ApiException{"Could not get message"}
}

// parseSingleMessageResponse parses the XML response for a single message and returns a structs.Message.
func parseSingleMessageResponse(body io.Reader) (structs.Message, error) {
	var response structs.ShowMessageResponse
	if err := xml.NewDecoder(body).Decode(&response); err != nil {
		return structs.Message{}, err
	}

	return response.Message, nil
}
