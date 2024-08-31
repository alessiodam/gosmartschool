package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func (client *SmartSchoolClient) wsOnError(c *websocket.Conn, err error) {
	client.websocketLogger.Error("Error:", err)
	if client.OnErrorHandler != nil {
		client.OnErrorHandler(c, err)
	}
}

func (client *SmartSchoolClient) wsOnClose(c *websocket.Conn, closeStatusCode int, closeMessage string) {
	client.websocketLogger.Info(fmt.Sprintf("WebSocket connection closed: %d - %s", closeStatusCode, closeMessage))
	if client.OnCloseHandler != nil {
		client.OnCloseHandler(c, closeStatusCode, closeMessage)
	}
}

func (client *SmartSchoolClient) wsOnMessage(c *websocket.Conn, message []byte) {
	client.websocketLogger.Debug("Received message:", string(message))

	var messageData map[string]interface{}
	if err := json.Unmarshal(message, &messageData); err != nil {
		client.websocketLogger.Error("Error unmarshalling message:", err)
		return
	}

	messageType, typeOk := messageData["type"].(string)
	messageRequest, requestOk := messageData["request"].(string)

	handled := false
	if typeOk && requestOk {
		switch messageType {
		case "auth":
			if messageRequest == "getToken" {
				client.websocketLogger.Info("Authentication successful!")
				handled = true
			}
		case "notificationListStart":
			client.websocketLogger.Info("Notification list started.")
			handled = true
		case "getNotificationConfig":
			configMessage := map[string]interface{}{
				"type":      "setConfig",
				"queueUuid": uuid.New().String(),
			}
			if err := c.WriteJSON(configMessage); err != nil {
				client.websocketLogger.Error("Error while sending config message:", err)
			}
			handled = true
		}
	}

	// If the message wasn't handled, invoke the user-defined handler
	if !handled && client.OnMessageHandler != nil {
		client.OnMessageHandler(c, messageData)
	}

	// Call the ReceivedMessageCallback if it's defined
	if client.ReceivedMessageCallback != nil {
		client.ReceivedMessageCallback(messageData)
	} else if text, ok := messageData["text"].(string); ok {
		client.websocketLogger.Info("Received message:", text)
	}
}

func (client *SmartSchoolClient) RunWebsocket() {
	client.websocketLogger.Info("Connecting to WebSocket")

	c, _, err := websocket.DefaultDialer.Dial("wss://nodejs-gs.smartschool.be/smsc/websocket", nil)
	if err != nil {
		client.websocketLogger.Fatal("Error while connecting to WebSocket:", err)
		return
	}
	defer func(c *websocket.Conn) {
		if err := c.Close(); err != nil {
			client.websocketLogger.Error("Error while closing WebSocket connection:", err)
		}
	}(c)

	//client.wsOnOpen(c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				client.wsOnClose(c, websocket.CloseNormalClosure, "Connection closed unexpectedly")
			} else {
				client.wsOnError(c, err)
			}
			break
		}

		client.wsOnMessage(c, message)
	}
}
