package obUtils

import (
	"fmt"
	"obAggregator/commonUtils"
	"obAggregator/protoMessages/messages"
)

func HandleDepthUpdates(exchange, instrument string) (*messages.OrderbookResponse, error) {
	switch exchange {
	case "BINANCE", "DERIBIT":
		//unmarshal to binance depth update struct
		//convert it to l2Update struct and return the value
		return generateOb(exchange, instrument), nil
	default:
		return nil, fmt.Errorf("exchange not implemented %s", exchange)
	}
}

func generateOb(exchange, instrument string) *messages.OrderbookResponse {
	asks := generateLevelBySide(instrument, "ask", 4)
	bids := generateLevelBySide(instrument, "bid", 4)
	ob := &messages.OrderbookResponse{
		Asks:         asks,
		Bids:         bids,
		ReceivedTime: commonUtils.CurrentTimestampMilli(),
		Instrument:   instrument,
		Exchange:     exchange,
	}
	return ob
}

func generateLevelBySide(instrument, side string, level int) []*messages.OrderbookLevel {
	var price float64
	if instrument == "BTC_USDT" {
		price = 28000
	} else if instrument == "ETH_USDT" {
		price = 1900
	}
	levels := []*messages.OrderbookLevel{}
	for i := 0; i < level; i++ {
		price = commonUtils.AddBpsBySide(price, 0.09, side)
		levels = append(levels, &messages.OrderbookLevel{Price: price, Size: 1})
	}
	return levels
}
