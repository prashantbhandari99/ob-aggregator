package ob

import (
	"fmt"
	"log"
	"obAggregator/actors/ob/obUtils"

	"github.com/asynkron/protoactor-go/actor"
)

// This actor is responsible for storing the latest orderbook state
// It also handles any requests from other actors for the orderbook information.
type OrderbookManager struct {
	actorName  string
	instrument string
	exchange   string
	ob         obUtils.Orderbook
	supervisor actor.SupervisorStrategy
	childWsPid *actor.PID
	stopChan   chan struct{}
}

func NewOrderbookManager(
	instrument,
	exchange string,
	supervisor actor.SupervisorStrategy,
) actor.Producer {
	return func() actor.Actor {
		return &OrderbookManager{
			exchange:   exchange,
			supervisor: supervisor,
			stopChan:   make(chan struct{}),
		}
	}
}

func (actr *OrderbookManager) init(context actor.Context) error {
	actr.actorName = context.Self().Id
	log.Println("=======started obm=======", actr.actorName)
	err := actr.startWsActor(context)
	if err != nil {
		return err
	}
	return nil
}

func (actr *OrderbookManager) startWsActor(context actor.Context) error {
	serviceName := fmt.Sprintf("WSS-%s-%s", actr.exchange, actr.instrument)
	props := actor.
		PropsFromProducer(NewOrderbookWsActor(actr.exchange, actr.instrument, actr.stopChan), actor.WithSupervisor(actr.supervisor))
	pid, err := context.SpawnNamed(props, serviceName)
	if err != nil {
		return fmt.Errorf("failed to start child-actor for bot %v", err)
	}
	actr.childWsPid = pid
	return nil
}

func (actr *OrderbookManager) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Printf("%s:Started, initializing...", actr.actorName)
		if err := actr.init(context); err != nil {
			log.Printf("%s - failed to initialize :- %s", actr.actorName, err)
			context.Stop(context.Self())
			return
		}
	case *actor.Stopping:
		log.Printf("%s-%s - Stopping...", actr.actorName, actr.instrument)
	case *actor.Stopped:
		log.Printf("%s - Stopped...", actr.actorName)
	case *actor.Restarting:
		log.Printf("%s - Restarting...", actr.actorName)
	case *actor.DeadLetterEvent:
		log.Printf("%s - DeadLetterEvent, %s", actr.actorName, fmt.Sprintf("UndeliveredMessage - %v", msg))
	case *actor.DeadLetterResponse:
		log.Printf("%s - DeadLetterResponse %s ", actr.actorName, fmt.Sprintf("UndeliveredMessage - %T", msg))
	case obUtils.Orderbook:
		actr.handleObUpdate(msg, context)
	case *obUtils.OrderbookRequest:
		actr.onObRequest(context)
	case *actor.Terminated:
		log.Printf("Terminated %v...", msg)
	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %T", msg))
	}
}

func (actr *OrderbookManager) handleObUpdate(msg obUtils.Orderbook,
	context actor.Context) {
	actr.ob = msg
}

func (actr *OrderbookManager) onObRequest(context actor.Context) {
	context.Respond(actr.ob)
}
