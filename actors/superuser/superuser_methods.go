package superuser

import (
	"fmt"
	"obAggregator/actors/user"
	"obAggregator/commonUtils"
	"obAggregator/protoMessages/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func (actr *SuperUser) startUsers(context actor.Context) error {
	for i := 0; i < actr.numUsers; i++ {
		props := actor.
			PropsFromProducer(user.NewUser(), actor.WithSupervisor(actr.supervisor))
		pid, err := context.SpawnNamed(props, fmt.Sprintf("ENDUSER-%v", i))
		if err != nil {
			return fmt.Errorf("failed to start enduser  error %v", err)
		}
		actr.childObPid = append(actr.childObPid, pid)
	}
	return nil
}

func (actr *SuperUser) onStop(context actor.Context) {
	commonUtils.PoisonActor(context, actr.childObPid)
	// context.Poison(actr.childWsPid)
}

func (actr *SuperUser) publishObToUser(context actor.Context, ob *messages.OrderbookResponse) {
	context.ActorSystem().EventStream.Publish(ob)
}
