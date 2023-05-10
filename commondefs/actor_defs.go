package commondefs

type ActorMessageType uint32

const (
	SWITCH_EXCHANGE ActorMessageType = iota
	OB_HEALTH
)

type ActorMessage struct {
	Type  ActorMessageType
	Value interface{}
}
