package ob

import (
	"fmt"
	"log"
	"obAggregator/actors/ob/obUtils"
	"obAggregator/protoMessages/messages"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gorilla/websocket"
)

// This actor processes updates received from a websocket and maintains a local order book.
// It also sends the latest update to the order book manager.
type OrderbookWsActor struct {
	conn                            *websocket.Conn
	exchange, instrument, actorName string
	// ob                              obUtils.Orderbook
	keepProcessing *atomic.Bool
}

func NewOrderbookWsActor(exchange, instrument string,
	keepProcessing *atomic.Bool,
) actor.Producer {
	return func() actor.Actor {
		return &OrderbookWsActor{
			exchange:       exchange,
			instrument:     instrument,
			keepProcessing: keepProcessing,
		}
	}
}

func (actr *OrderbookWsActor) init(context actor.Context) error {
	actr.actorName = context.Self().Id
	log.Printf("=========Started ob ws actor %s========", actr.actorName)
	actr.connectWs(context)
	return nil
}

func (actr *OrderbookWsActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		if err := actr.init(context); err != nil {
			log.Printf("%s - failed to initialize :- %s", actr.actorName, err)
			context.Stop(context.Self())
			return
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
	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %T", msg))
	}
}

func (actr *OrderbookWsActor) onStop(context actor.Context) {
	if actr.conn != nil {
		actr.conn.Close()
	}
}

func (actr *OrderbookWsActor) connectWs(context actor.Context) error {
	//init actr.ob
	// dialer := &websocket.Dialer{}
	// conn, _, err := dialer.Dial("endpoint", nil)
	// if err != nil {
	// 	return err
	// }
	// actr.conn = conn
	actr.readMessages(context)
	return nil
}

func (actr *OrderbookWsActor) handleDisconnection(context actor.Context) {
	log.Println("Disconnection ws ", actr.actorName)
	actr.onStop(context)
	time.Sleep(5 * time.Second)
	actr.connectWs(context)
}

func (actr *OrderbookWsActor) readMessages(context actor.Context) {
	ticker := time.NewTicker(400 * time.Millisecond) // send a ping every 10 seconds
	defer ticker.Stop()
	actr.keepProcessing.Store(true)
	for actr.keepProcessing.Load() {

		select {
		case <-ticker.C:
			err := actr.handleL2Update(context)
			if err != nil {
				log.Println(err)
				return
			}
		default:
			// 	actr.writeMessage([]byte("ping"))
			// default:
			// 	_, message, err := actr.conn.ReadMessage()
			// 	if err != nil {
			// 		return
		}

		// }
	}
}

func (actr *OrderbookWsActor) handleL2Update(context actor.Context) error {
	ob, err := obUtils.HandleDepthUpdates(actr.exchange, actr.instrument)
	if err != nil {
		return err
	}
	// apply l2Updates to the local orderbook (actr.ob). If any event drop is
	// detected, the program will disconnect and reconnect the websocket using actr.handleDisconnecion. This
	// will also reset the state of the local orderbook (actr.ob).
	actr.sendObUpdate(context, ob)
	return nil
}

func (actr *OrderbookWsActor) sendObUpdate(context actor.Context, ob *messages.OrderbookResponse) {
	if len(ob.Asks) > 0 && len(ob.Bids) > 0 {
		context.Send(context.Parent(), ob) // sends latest orderbook to parent orderbook manager
	}
}

func (actr *OrderbookWsActor) writeMessage(message []byte) {
	err := actr.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("failed to write message", actr.actorName, string(message))
	}
}
