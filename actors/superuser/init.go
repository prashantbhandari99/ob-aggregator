package superuser

import (
	"log"
	"os"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
)

// this actor receives ob from delta actor and publishes latest ob to all users
type SuperUser struct {
	actorName  string
	supervisor actor.SupervisorStrategy
	childObPid []*actor.PID
	numUsers   int
}

func NewDSuperUser(
	supervisor actor.SupervisorStrategy,
) actor.Producer {
	return func() actor.Actor {
		numUsers, err := strconv.Atoi(os.Getenv("NUM_USERS"))
		if err != nil {
			numUsers = 10
		}
		return &SuperUser{
			supervisor: supervisor,
			numUsers:   numUsers,
		}
	}
}

func (actr *SuperUser) init(context actor.Context) error {
	actr.actorName = context.Self().Id
	log.Println("=======started superuser=======",
		actr.actorName)
	actr.startUsers(context)
	return nil
}
