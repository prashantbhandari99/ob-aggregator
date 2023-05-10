package user

import (
	"fmt"
	"log"
	"obAggregator/protoMessages/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func (actr *User) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		if err := actr.init(context); err != nil {
			log.Printf("%s - failed to initialize :- %s", actr.actorName, err)
			context.Stop(context.Self())
			return
		}
	case *actor.Stopping:
		actr.onStop(context)
		log.Printf("%s - Stopping...", actr.actorName)
	case *actor.Stopped:
		log.Printf("%s - Stopped...", actr.actorName)
	case *actor.Restarting:
		log.Printf("%s - Restarting...", actr.actorName)
	case *actor.DeadLetterEvent:
		log.Printf("%s - DeadLetterEvent, %s", actr.actorName, fmt.Sprintf("UndeliveredMessage - %v", msg))
	case *actor.DeadLetterResponse:
		log.Printf("%s - DeadLetterResponse %s ", actr.actorName, fmt.Sprintf("UndeliveredMessage - %T", msg))
	case *actor.Terminated:
		log.Printf("Terminated %v...", msg)
	case *messages.OrderbookResponse:

	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %T", msg))
	}
}
