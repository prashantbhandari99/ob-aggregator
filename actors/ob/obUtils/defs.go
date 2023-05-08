package obUtils

type ArbOrderType uint32

const (
	BINANCE ArbOrderType = iota
)

type Orderbook struct {
	Asks                     OrderBookLevel
	Bids                     OrderBookLevel
	SeqNumber, PrevSeqNumber int64
	RecvTime                 int64
}

type OrderBookLevel struct {
	Price float64
	Size  float64
}
