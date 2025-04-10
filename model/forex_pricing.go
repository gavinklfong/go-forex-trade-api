package model

type ForexPricing struct {
	BaseCurrency    string
	CounterCurrency string
	BuyPip          float32
	SellPip         float32
}

func (p *ForexPricing) GetSpread() float32 {
	return p.BuyPip - p.SellPip
}
