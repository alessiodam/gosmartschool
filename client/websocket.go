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
	client.ApiLogger.Info("Requesting token from API")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/Topnav/Node/getToken", client.domain), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", client.PhpSessId))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.ApiLogger.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		token, _ := io.ReadAll(resp.Body)
		client.userToken = string(token)
		client.ApiLogger.Info("Websocket token received")
		return client.userToken, nil
	}

	client.ApiLogger.Error("Could not get websocket token")
	return "", &ApiException{"Could not get websocket token"}
}

func (client *SmartSchoolClient) wsOnOpen(_ *websocket.Conn) error {
	client.WebsocketLogger.Info("WebSocket connection opened")

	apiToken, err := client.GetWebsocketTokenFromAPI()
	if err != nil {
		return err
	}
	client.websocketToken = apiToken
	return nil
}

func (client *SmartSchoolClient) wsOnError(c *websocket.Conn, err error) {
	client.WebsocketLogger.Error("Error:", err)
	if client.wsOnErrorHandler != nil {
		client.wsOnErrorHandler(c, err)
	}
}

func (client *SmartSchoolClient) wsOnClose(c *websocket.Conn, closeStatusCode int, closeMessage string) {
	client.WebsocketLogger.Info(fmt.Sprintf("WebSocket connection closed: %d - %s", closeStatusCode, closeMessage))
	if client.wsOnCloseHandler != nil {
		client.wsOnCloseHandler(c, closeStatusCode, closeMessage)
	}
}

func (client *SmartSchoolClient) wsOnMessage(c *websocket.Conn, message []byte) {
	client.WebsocketLogger.Debug("Received message:", string(message))

	var messageData map[string]interface{}
	if err := json.Unmarshal(message, &messageData); err != nil {
		client.WebsocketLogger.Error("Error unmarshalling message:", err)
		return
	}

	messageType, typeOk := messageData["type"].(string)
	messageRequest, requestOk := messageData["request"].(string)

	handled := false
	if typeOk && requestOk {
		switch messageType {
		case "auth":
			if messageRequest == "getToken" {
				client.WebsocketLogger.Info("Authentication successful!")
				handled = true
			}
		case "notificationListStart":
			client.WebsocketLogger.Info("Notification list started.")
			handled = true
		case "getNotificationConfig":
			configMessage := map[string]interface{}{
				"type":      "setConfig",
				"queueUuid": uuid.New().String(),
			}
			if err := c.WriteJSON(configMessage); err != nil {
				client.WebsocketLogger.Error("Error while sending config message:", err)
			}
			handled = true
		}
	}

	// If the message wasn't handled, invoke the user-defined handler
	if !handled && client.wsOnMessageHandler != nil {
		client.wsOnMessageHandler(c, messageData)
	}

	// Call the ReceivedMessageCallback if it's defined
	if client.ReceivedMessageCallback != nil {
		client.ReceivedMessageCallback(messageData)
	} else if text, ok := messageData["text"].(string); ok {
		client.WebsocketLogger.Info("Received message:", text)
	}
}

func (client *SmartSchoolClient) RunWebsocket() {
	client.WebsocketLogger.Info("Connecting to WebSocket")

	wsConn, _, err := websocket.DefaultDialer.Dial("wss://nodejs-gs.smartschool.be/smsc/websocket", nil)
	if err != nil {
		client.WebsocketLogger.Fatal("Error while connecting to WebSocket:", err)
		return
	}
	defer func(c *websocket.Conn) {
		if err := c.Close(); err != nil {
			client.WebsocketLogger.Error("Error while closing WebSocket connection:", err)
		}
	}(wsConn)

	err = client.wsOnOpen(wsConn)
	if err != nil {
		return
	}

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				client.wsOnClose(wsConn, websocket.CloseNormalClosure, "Connection closed unexpectedly")
			} else {
				client.wsOnError(wsConn, err)
			}
			break
		}

		client.wsOnMessage(wsConn, message)
	}
}
