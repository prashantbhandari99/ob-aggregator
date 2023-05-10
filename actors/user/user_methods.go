package user

import (
	"log"
	"obAggregator/protoMessages/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func (actr *User) onStop(context actor.Context) {
	context.ActorSystem().EventStream.Unsubscribe(actr.handler)
}

func (actr *User) subOb(context actor.Context) {
	handler := func(message interface{}) {
		obUpdate := message.(*messages.OrderbookResponse)
		log.Printf("Recieved orderbook update %s, %s; best ask %.5f--best bid %.5f %s",
			obUpdate.Exchange, obUpdate.Instrument, obUpdate.Asks[0].Price, obUpdate.Bids[0].Price, actr.actorName)
	}
	actr.handler = context.ActorSystem().EventStream.Subscribe(handler)
}
