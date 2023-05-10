package delta

import (
	"fmt"
	"obAggregator/actors/ob"
	"obAggregator/commonUtils"
	"obAggregator/protoMessages/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func (actr *DeltaWorker) initObm(context actor.Context) error {
	for _, exchange := range actr.exchanges {
		for _, instrument := range actr.instrumentsToSub {
			pid, err := actr.startObm(context, instrument, exchange)
			if err != nil {
				return err
			}
			actr.childObPid = append(actr.childObPid, pid)
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func (actr *DeltaWorker) startObm(context actor.Context, instrument, exchange string) (*actor.PID, error) {
	props := actor.
		PropsFromProducer(ob.NewOrderbookManager(exchange, instrument, actr.supervisor), actor.WithSupervisor(actr.supervisor))
	pid, err := context.SpawnNamed(props, fmt.Sprintf("OBM-%s-%s", exchange, instrument))
	if err != nil {
		return nil, fmt.Errorf("failed to start obm %s, %s error %v", exchange, instrument, err)
	}
	return pid, nil
}

func (actr *DeltaWorker) onStop(context actor.Context) {
	commonUtils.PoisonActor(context, actr.childObPid)
	commonUtils.CancelSchedulers(actr.cancelSchedulers)
	commonUtils.Unsubscribe(context, actr.subHandlers)
}

func (actr *DeltaWorker) subObs(context actor.Context) {
	for _, exchange := range actr.exchanges {
		for _, instrument := range actr.instrumentsToSub {
			actr.subOb(context, exchange, instrument)
		}
	}
}

func (actr *DeltaWorker) subOb(context actor.Context, exchange, instrument string) {
	handler := func(message interface{}) {
		obUpdate := message.(*messages.OrderbookResponse)
		if instrument == obUpdate.Instrument && exchange == obUpdate.Exchange {
			actr.publishObToUser(context, obUpdate)
		}
	}
	predicate := func(event interface{}) bool {
		_, ok := event.(*messages.OrderbookResponse)
		return ok
	}
	actr.subHandlers = append(actr.subHandlers,
		context.ActorSystem().EventStream.SubscribeWithPredicate(handler, predicate))
}

func (actr *DeltaWorker) publishObToUser(context actor.Context, ob *messages.OrderbookResponse) {
	if ob.Exchange == actr.currentExchange {
		context.Send(actr.superuserPid, ob)
	}
}
