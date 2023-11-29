package binanceusdm

import (
	"time"

	"github.com/shopspring/decimal"
)

type Spread struct {
	WsReceived time.Time
	Time       decimal.Decimal
	BidVolume  decimal.Decimal
	Bid        decimal.Decimal
	Ask        decimal.Decimal
	AskVolume  decimal.Decimal
}
