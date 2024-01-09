package websocket

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/util/log"
)

func HandleConnectionSwitch(ct context.Context, conn *webSocketConn) {
	// * Extract player id from context
	playerId := ct.Value("playerId").(primitive.ObjectID)

	// * Connection switch
	if conn.Conn != nil {
		log.Debug("Connection switch", "playerId", playerId)
		conn.Emit(&OutboundMessage{
			Event:   ConnectionSwitchEvent,
			Payload: nil,
		})

		conn.Mutex.Lock()
		if err := conn.Conn.Close(); err != nil {
			log.Error("Unable to close connection", err)
		}
	}
}
