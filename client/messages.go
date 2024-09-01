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

// parseMessageListResponse parses the XML response for a list of messages and returns a slice of structs.Message.
func (client *SmartSchoolClient) parseMessageListResponse(body io.Reader) ([]structs.Message, error) {
	var server structs.XMLServer

	if err := xml.NewDecoder(body).Decode(&server); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	var messages []structs.Message
	for _, action := range server.Response.Actions {
		messages = append(messages, action.Data.Messages...)
	}

	return messages, nil
}

func (client *SmartSchoolClient) constructListMessagesCommand(boxType string, inboxId int) string {
	return fmt.Sprintf(`
		<request>
			<command>
				<subsystem>postboxes</subsystem>
				<action>message list</action>
				<params>
					<param name="boxType"><![CDATA[%s]]></param>
					<param name="boxID"><![CDATA[%d]]></param>
					<param name="sortField"><![CDATA[date]]></param>
					<param name="sortKey"><![CDATA[desc]]></param>
					<param name="poll"><![CDATA[false]]></param>
					<param name="poll_ids"><![CDATA[]]></param>
					<param name="layout"><![CDATA[new]]></param>
				</params>
			</command>
		</request>
	`, boxType, inboxId)
}

// constructGetMessageCommand constructs the XML command to retrieve a message by its ID.
func (client *SmartSchoolClient) constructGetMessageCommand(messageID string) string {
	return fmt.Sprintf(`
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
		</request>`, messageID)
}

// constructDeleteMessageCommand constructs the XML command to delete a message by its ID.
func (client *SmartSchoolClient) constructDeleteMessageCommand(messageID string) string {
	return fmt.Sprintf(`
		<request>
			<command>
				<subsystem>postboxes</subsystem>
				<action>quick delete</action>
				<params>
					<param name="msgID"><![CDATA[%s]]></param>
				</params>
			</command>
		</request>`, messageID)
}

// ListMessages retrieves a list of messages from the API and returns a slice of structs.Message.
//
// Returns last 50 messages because of SmartSchool API limitations.
func (client *SmartSchoolClient) ListMessages(boxType string, inboxId int) ([]structs.Message, error) {
	client.ApiLogger.Info("Requesting messages from API")

	data := url.Values{}
	data.Set("command", client.constructListMessagesCommand(boxType, inboxId))

	resp, body, err := client.sendXmlRequest("POST", "/?module=Messages&file=dispatcher", data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.ApiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get messages, status: %d", resp.StatusCode)
	}

	if body == "" {
		return nil, fmt.Errorf("empty response body")
	}

	parsedMessages, err := client.parseMessageListResponse(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	client.ApiLogger.Info("Messages received")
	return parsedMessages, nil
}

// DeleteMessageByID deletes a single message by its ID.
func (client *SmartSchoolClient) DeleteMessageByID(messageID string) error {
	client.ApiLogger.Info("Requesting message deletion from API")

	data := url.Values{}
	data.Set("command", client.constructDeleteMessageCommand(messageID))

	resp, _, err := client.sendXmlRequest("POST", "/?module=Messages&file=dispatcher", data.Encode(), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.ApiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not delete message, status: %d", resp.StatusCode)
	}

	client.ApiLogger.Info("Message deleted")
	return nil
}
