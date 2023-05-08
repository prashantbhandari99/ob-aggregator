package obUtils

import "fmt"

func HandleDepthUpdates(update []byte, exchange string) (*Orderbook, error) {
	switch exchange {
	case "BINANCE":
		//unmarshal to binance depth update struct
		//convert it to l2Update struct and return the value
		return &Orderbook{}, nil
	default:
		return nil, fmt.Errorf("exchange not implemented %s", exchange)
	}
}
