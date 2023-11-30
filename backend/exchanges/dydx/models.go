package dydx

import "github.com/shopspring/decimal"

type InitResp struct {
	Channel      string       `json:"channel"`
	ConnectionId string       `json:"connection_id"`
	Contents     InitContents `json:"contents"`
	Id           string       `json:"id"`
	Type         string       `json:"type"`
}

type InitContents struct {
	Asks []InitAsk `json:"asks"`
	Bids []InitBid `json:"bids"`
}

type InitAsk struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type InitBid struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type UpdateResp struct {
	Channel      string         `json:"channel"`
	ConnectionId string         `json:"connection_id"`
	Contents     UpdateContents `json:"contents"`
	Id           string         `json:"id"`
	Type         string         `json:"type"`
}

type UpdateContents struct {
	Asks   []interface{} `json:"asks"`
	Bids   []interface{} `json:"bids"`
	Offset string        `json:"offset"`
}

type Book struct {
	Asks []Ask
	Bids []Bid
}

type Bid struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

type Ask struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

func (a *Bid) Equal(b Bid) bool {
	if a.Price.Equal(b.Price) && a.Size.Equal(b.Size) {
		return true
	} else {
		return false
	}
}

func (a *Ask) Equal(b Ask) bool {
	if a.Price.Equal(b.Price) && a.Size.Equal(b.Size) {
		return true
	} else {
		return false
	}
}
