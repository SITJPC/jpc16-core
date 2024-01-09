package websocket

import (
	"fmt"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type webSocketInstance struct {
	Players map[primitive.ObjectID]*webSocketConn
}

type webSocketConn struct {
	Conn  *websocket.Conn
	Mutex *sync.Mutex
}

func (r *webSocketConn) Emit(payload *OutboundMessage) {
	if r.Conn == nil || r.Conn.Conn == nil {
		return
	}

	r.Mutex.Lock()
	if err := r.Conn.WriteJSON(payload); err != nil {
		logrus.Warn(fmt.Sprintf("WRITING MESSAGE FAILURE: %s", err.Error()))
	}
	r.Mutex.Unlock()
}

func Emit(playerId primitive.ObjectID, payload *OutboundMessage) {
	conn := Instance.Players[playerId]
	if conn == nil {
		return
	}

	conn.Emit(payload)
}

var Instance *webSocketInstance

func Init() {
	Instance = &webSocketInstance{
		Players: make(map[primitive.ObjectID]*webSocketConn),
	}
}
