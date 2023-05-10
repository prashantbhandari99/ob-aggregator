package commonUtils

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
	"github.com/asynkron/protoactor-go/scheduler"
)

func CancelSchedulers(cancels []scheduler.CancelFunc) {
	for _, cancel := range cancels {
		cancel()
	}
}

func PoisonActor(context actor.Context, pids []*actor.PID) {
	for _, pid := range pids {
		context.Poison(pid)
	}
}

func Unsubscribe(context actor.Context, handlers []*eventstream.Subscription) {
	for _, handler := range handlers {
		context.ActorSystem().EventStream.Unsubscribe(handler)
	}
}
