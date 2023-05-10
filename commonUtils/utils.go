package commonUtils

func AddBpsBySide(price, bps float64, side string) float64 {
	if side == "ask" {
		return price + (price * bps / 10000)
	}
	return price - (price * bps / 10000)
}
