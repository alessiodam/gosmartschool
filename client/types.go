package client

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ApiException struct {
	Message string
}

func (e *ApiException) Error() string {
	return e.Message
}

type AuthException struct {
	ApiException
}

type WebSocketErrorHandler func(c *websocket.Conn, err error)
type WebSocketCloseHandler func(c *websocket.Conn, closeStatusCode int, closeMessage string)
type WebSocketMessageHandler func(c *websocket.Conn, message map[string]interface{})

type SmartSchoolUser struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Class     string
}

type SmartSchoolClient struct {
	domain     string
	platformId string

	UniqueUsc string
	PhpSessId string
	userToken string

	User SmartSchoolUser

	WriteApiLogs    bool
	ApiLogger       *logrus.Logger
	WebsocketLogger *logrus.Logger
	AuthLogger      *logrus.Logger

	websocketToken          string
	ReceivedMessageCallback func(message map[string]interface{})
	wsOnErrorHandler        WebSocketErrorHandler
	wsOnCloseHandler        WebSocketCloseHandler
	wsOnMessageHandler      WebSocketMessageHandler
}
