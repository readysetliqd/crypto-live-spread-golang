package data

import (
	"time"

	"github.com/shopspring/decimal"
)

type Spread struct {
	WsReceived time.Time
	Time       decimal.Decimal // seconds since epoch
	Bid        decimal.Decimal
	Ask        decimal.Decimal
	BidVolume  decimal.Decimal
	AskVolume  decimal.Decimal
}
