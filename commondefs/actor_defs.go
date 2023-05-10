package commondefs

// ActorMessageType actor message type
type ActorMessageType uint32

// actor message type
const (
	SWITCH_EXCHANGE ActorMessageType = iota
)

type ActorMessage struct {
	Type  ActorMessageType
	Value interface{}
}
