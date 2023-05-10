package delta

import (
	"fmt"
	"log"
	"obAggregator/commondefs"

	"github.com/asynkron/protoactor-go/actor"
)

func (actr *DeltaWorker) Receive(context actor.Context) {
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
	case *commondefs.ActorMessage:
		actr.onMessage(msg)
	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %T", msg))
	}
}

func (actr *DeltaWorker) onMessage(msg *commondefs.ActorMessage) {
	switch msg.Type {
	case commondefs.SWITCH_EXCHANGE:
		for _, exchange := range actr.exchanges {
			if actr.currentExchange != exchange {
				actr.currentExchange = exchange
				break
			}
		}
	default:
		log.Println("msg type not implemented", actr.actorName, msg)
	}
}
