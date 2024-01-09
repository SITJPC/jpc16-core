package websocket

import (
	"context"
	"sync"

	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	playerRepo "jpc16-core/repository/player"
	"jpc16-core/util/log"
	"jpc16-core/util/value"
)

var InitialFunc func(ct context.Context)

func Serve(conn *websocket.Conn) {
	// * Construct context
	ct := context.Background()

	// * Parse connection
	playerId := conn.Query("playerId")
	playerIdObj, err := primitive.ObjectIDFromHex(playerId)
	if err != nil {
		log.Error("Invalid player ID", err)
		return
	}
	ct = context.WithValue(ct, "playerId", playerIdObj)

	// * Validate player
	player, errr := playerRepo.FindById(&playerIdObj)
	if errr != nil {
		log.Error("Unable to find player", err)
		return
	}
	ct = context.WithValue(ct, "player", player)

	// * Handle connection switch
	cnn := Instance.Players[playerIdObj]
	if cnn != nil {
		HandleConnectionSwitch(ct, cnn)
	} else {
		cnn = &webSocketConn{
			Conn:  nil,
			Mutex: new(sync.Mutex),
		}
		Instance.Players[playerIdObj] = cnn
	}

	// * Assign connection
	cnn.Mutex.Lock()
	cnn.Conn = conn
	cnn.Mutex.Unlock()

	// * Emit initial message
	InitialFunc(ct)

	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if t != websocket.TextMessage {
			break
		}

		cnn.Emit(&OutboundMessage{
			Event:   EchoEvent,
			Payload: p,
		})
	}

	// * Close connection
	if err := conn.Close(); err != nil {
		log.Error("Unable to close connection", err)
	}

	// * Reset player connection
	cnn.Conn = nil

	// * Unlock in case of connection switch
	if value.MutexLocked(cnn.Mutex) {
		cnn.Mutex.Unlock()
	}
}
