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
	domain          string
	platformId      string
	PhpSessId       string
	Pid             string
	userID          string
	userToken       string
	apiLogger       *logrus.Logger
	websocketLogger *logrus.Logger
	authLogger      *logrus.Logger

	ReceivedMessageCallback func(message map[string]interface{})
	OnErrorHandler          WebSocketErrorHandler
	OnCloseHandler          WebSocketCloseHandler
	OnMessageHandler        WebSocketMessageHandler
}
