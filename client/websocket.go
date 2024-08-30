package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

func (client *SmartSchoolClient) GetWebsocketTokenFromAPI() (string, error) {
	client.apiLogger.Info("Requesting token from API")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/Topnav/Node/getToken", client.domain), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s; pid=%s", client.PhpSessId, client.Pid))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.apiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		token, _ := io.ReadAll(resp.Body)
		client.userToken = string(token)
		client.apiLogger.Info("Websocket token received")
		return client.userToken, nil
	}

	client.apiLogger.Error("Could not get websocket token")
	return "", &ApiException{"Could not get websocket token"}
}

func (client *SmartSchoolClient) RunWebsocket() {
	client.websocketLogger.Info("Connecting to websocket ")

	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://nodejs-gs.smartschool.be/smsc/websocket"), nil)
	if err != nil {
		client.websocketLogger.Fatal("Error while connecting to WebSocket:", err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			client.websocketLogger.Error(err)
		}
	}(c)

	authMessage := map[string]interface{}{
		"type":    "auth",
		"request": "checkToken",
		"token":   client.userToken,
	}

	if err := c.WriteJSON(authMessage); err != nil {
		client.websocketLogger.Fatal("Error while sending auth message:", err)
	}

	for {
		var message map[string]interface{}
		err := c.ReadJSON(&message)
		if err != nil {
			client.websocketLogger.Fatal("Error while reading message:", err)
		}

		client.websocketLogger.Debug("Received message:", message)
		client.handleWebSocketMessage(message)
	}
}

func (client *SmartSchoolClient) handleWebSocketMessage(message map[string]interface{}) {
	messageType := message["type"].(string)
	messageRequest := message["request"].(string)

	if messageType == "auth" && messageRequest == "getToken" {
		client.websocketLogger.Info("Authentication successful!")
	} else if messageType == "notificationListStart" {
		client.websocketLogger.Info("Notification list started.")
	} else if messageType == "getNotificationConfig" {
		// configMessage := map[string]interface{}{
		// 	"type":      "setConfig",
		// 	"queueUuid": uuid.New().String(),
		// }
		// client.runWebSocket().WriteJSON(configMessage)
	}

	if client.receivedMessageCallback != nil {
		client.receivedMessageCallback(message)
	} else {
		client.websocketLogger.Info("Received message:", message["text"].(string))
	}
}
