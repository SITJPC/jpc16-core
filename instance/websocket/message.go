package websocket

type InboundEvent string
type OutboundEvent string

const (
	ConnectionSwitchEvent OutboundEvent = "general/switch"
	EchoEvent             OutboundEvent = "general/echo"
	PlayerStateEvent      OutboundEvent = "player/state"
)

type OutboundMessage struct {
	Event   OutboundEvent `json:"event"`
	Payload any           `json:"payload"`
}
