package delta

import (
	"log"
	"obAggregator/commondefs"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/asynkron/protoactor-go/scheduler"
)

//Delta actor subscribes to orderbook manager
//whenever update is recieved it sends the message to user

type DeltaWorker struct {
	actorName, currentExchange  string
	supervisor                  actor.SupervisorStrategy
	childObPid                  []*actor.PID
	instrumentsToSub, exchanges []string
	superuserPid                *actor.PID
	cancelSchedulers            []scheduler.CancelFunc
	subHandlers                 []*eventstream.Subscription
}

func NewDeltaWorker(
	instruments,
	exchanges string,
	supervisor actor.SupervisorStrategy,
) actor.Producer {
	return func() actor.Actor {
		instrumentsToSub := strings.Split(instruments, ",")
		exchanges := strings.Split(exchanges, ",")
		return &DeltaWorker{
			exchanges:        exchanges,
			supervisor:       supervisor,
			instrumentsToSub: instrumentsToSub,
			currentExchange:  exchanges[0],
		}
	}
}

func (actr *DeltaWorker) init(context actor.Context) error {
	actr.superuserPid = &actor.PID{
		Address: "0.0.0.0:9082",
		Id:      "SUPERUSER",
	}
	actr.actorName = context.Self().Id
	log.Println("=======started delta=======",
		actr.actorName)
	err := actr.initObm(context)
	if err != nil {
		return err
	}
	timer := scheduler.NewTimerScheduler(context)
	actr.cancelSchedulers = append(actr.cancelSchedulers, timer.SendRepeatedly(time.Second*15, time.Second*15,
		context.Self(), &commondefs.ActorMessage{
			Type: commondefs.SWITCH_EXCHANGE,
		}))
	time.Sleep(5 * time.Second)
	actr.subObs(context)
	return nil
}
