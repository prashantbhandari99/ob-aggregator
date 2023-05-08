package ob

import (
	"fmt"
	"log"
	"obAggregator/actors/ob/obUtils"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/gorilla/websocket"
)

type OrderBookManager struct {
	actorName  string
	instrument string
	exchange   string
	// ob             obUtils.Orderbook
	keepProcessing atomic.Value
}

func NewOrderBookManager(
	instrument,
	exchange string,
) actor.Producer {
	return func() actor.Actor {
		return &OrderBookManager{
			actorName: fmt.Sprintf("%s-%s", exchange, instrument),
			exchange:  exchange,
		}
	}
}

func (actr *OrderBookManager) init(context actor.Context) error {
	timerScheduler := scheduler.NewTimerScheduler(context)
	timerScheduler.SendRepeatedly(
		time.Duration(40)*time.Second, time.Duration(5)*time.Second, context.Self(), &obUtils.OrderBookLevel{})
	actr.initOb()
	return nil
}

func (actr *OrderBookManager) initOb() {
	//initialize local orderbook
	go actr.startEventProcessing()
}

func (actr *OrderBookManager) stopEventProcessing() {
	actr.keepProcessing.Store(false)
}

func (actr *OrderBookManager) startEventProcessing() {
	// define the WebSocket endpoint
	done := make(chan struct{})
	defer close(done)
	endpoint := "wss://stream.binance.com:9443/ws/btcusdt@depth"

	// create a new WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatal("error connecting to WebSocket:", err)
	}
	defer conn.Close()
	// read messages from the WebSocket
	actr.keepProcessing.Store(true)
	for actr.keepProcessing.Load().(bool) {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			return
		}
		actr.handleL2Update(message)
	}
}

func (actr *OrderBookManager) handleL2Update(update []byte) error {
	_, err := obUtils.HandleDepthUpdates(update, actr.exchange)
	if err != nil {
		return err
	}
	// apply l2Updates to the local orderbook (actr.ob). If any event drop is
	// detected, the program will disconnect and reconnect the websocket. This
	// will also reset the state of the local orderbook (actr.ob).
	return nil
}

func (actr *OrderBookManager) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Printf("%s:Started, initializing...", actr.actorName)
		if err := actr.init(context); err != nil {
			log.Printf("%s - failed to initialize :- %s", actr.actorName, err)
			context.Stop(context.Self())
			return
		}
	case *obUtils.OrderBookLevel:
		log.Println("Haye oye")
		actr.stopEventProcessing()
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
	default:
		log.Printf("%s - MessageError : %s ", actr.actorName, fmt.Sprintf("message type not implemeneted - %v", msg))
	}
}
