package ob

import (
	"fmt"
	"log"
	"obAggregator/commondefs"
	"obAggregator/protoMessages/messages"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
)

// This actor is responsible for storing the latest orderbook state
// It also handles any requests from other actors for the orderbook information.
type OrderbookManager struct {
	actorName      string
	instrument     string
	exchange       string
	supervisor     actor.SupervisorStrategy
	childWsPid     *actor.PID
	keepProcessing *atomic.Bool
}

func NewOrderbookManager(
	exchange,
	instrument string,
	supervisor actor.SupervisorStrategy,
) actor.Producer {
	return func() actor.Actor {
		return &OrderbookManager{
			exchange:   exchange,
			instrument: instrument,
			supervisor: supervisor,
		}
	}
}

func (actr *OrderbookManager) init(context actor.Context) error {
	actr.actorName = context.Self().Id
	log.Println("=======started obm=======",
		actr.actorName)
	err := actr.startWsActor(context)
	if err != nil {
		return err
	}
	timer := scheduler.NewTimerScheduler(context)
	timer.SendRepeatedly(time.Hour*23, time.Hour*23, context.Self(), &commondefs.ActorMessage{
		Type: commondefs.OB_HEALTH,
	})
	return nil
}

func (actr *OrderbookManager) startWsActor(context actor.Context) error {
	actr.keepProcessing = &atomic.Bool{}
	serviceName := fmt.Sprintf("WSS-%s-%s", actr.exchange, actr.instrument)
	props := actor.
		PropsFromProducer(NewOrderbookWsActor(actr.exchange, actr.instrument, actr.keepProcessing),
			actor.WithSupervisor(actr.supervisor))
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
		if err := actr.init(context); err != nil {
			log.Printf("%s - failed to initialize :- %s", actr.actorName, err)
			context.Stop(context.Self())
			return
		}
	case *commondefs.ActorMessage:
		if msg.Type == commondefs.OB_HEALTH {
			actr.keepProcessing.Store(false)
		}
	case *actor.Stopping:
		actr.onStop(context)
		log.Printf("%s-%s - Stopping...", actr.actorName, actr.instrument)
	case *actor.Stopped:
		log.Printf("%s - Stopped...", actr.actorName)
	case *actor.Restarting:
		log.Printf("%s - Restarting...", actr.actorName)
	case *actor.DeadLetterEvent:
		log.Printf("%s - DeadLetterEvent, %s", actr.actorName, fmt.Sprintf("UndeliveredMessage - %v", msg))
	case *actor.DeadLetterResponse:
		log.Printf("%s - DeadLetterResponse %s ", actr.actorName, fmt.Sprintf("UndeliveredMessage - %T", msg))
	case *messages.OrderbookResponse:
		actr.handleObUpdate(msg, context)
	case *actor.Terminated:
		log.Printf("Terminated %v...", msg)
	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %T", msg))
	}
}

func (actr *OrderbookManager) handleObUpdate(msg *messages.OrderbookResponse,
	context actor.Context) {
	context.ActorSystem().EventStream.Publish(msg)
}

func (actr *OrderbookManager) onStop(context actor.Context) {
	actr.keepProcessing.Store(false)
	context.Poison(actr.childWsPid)
}
