package user

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
)

// receives latest book
type User struct {
	actorName string
	handler   *eventstream.Subscription
}

func NewUser() actor.Producer {
	return func() actor.Actor {
		return &User{}
	}
}

func (actr *User) init(context actor.Context) error {
	actr.actorName = context.Self().Id
	log.Println("=======started superuser=======",
		actr.actorName)
	actr.subOb(context)
	return nil
}
