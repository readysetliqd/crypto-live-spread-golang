package dydx

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

// Accepts a channel and a pair. Connects to websocket api to updated channel
// with the live bid/ask spread for that pair
func GetSpread(updateChannel chan data.Spread, pair string) {
	defer close(updateChannel)
	var initResp map[string]interface{}
	var err error
	c, err := dial(initResp)
	if err != nil {
		log.Fatal("Failed initial connect | ", err)
	}
	defer c.Close()

	payload := `
	{
		"type": "subscribe",
		"channel": "v3_orderbook",
		"id": "` + pair + `"
	}
	`
	log.Println(payload)
	initBookMsg := subscribeChannel(c, initResp, payload)
	// var resp Resp
	// json.Unmarshal([]byte(init), &resp)
	var initResponse InitResp
	err = json.Unmarshal(initBookMsg, &initResponse)
	if err != nil {
		log.Fatal(err)
	}
	initBook := initResponse.Contents
	initAsks := initBook.Asks
	initBids := initBook.Bids

	// build initial state of book
	var book Book
	for _, ask := range initAsks {
		price, err := decimal.NewFromString(ask.Price)
		if err != nil {
			log.Fatal("decimal.NewFromString error | ", err)
		}
		size, err := decimal.NewFromString(ask.Size)
		if err != nil {
			log.Fatal("decimal.NewFromString error | ", err)
		}
		newAsk := Ask{price, size}
		book.Asks = append(book.Asks, newAsk)
	}
	for _, bid := range initBids {
		price, err := decimal.NewFromString(bid.Price)
		if err != nil {
			log.Fatal("decimal.NewFromString error | ", err)
		}
		size, err := decimal.NewFromString(bid.Size)
		if err != nil {
			log.Fatal("decimal.NewFromString error | ", err)
		}
		newBid := Bid{price, size}
		book.Bids = append(book.Bids, newBid)
	}

	// build initial bestBid and bestAsk from initial state of book
	var bestAsk Ask
	var bestBid Bid
	bestAsk.Price = book.Asks[0].Price
	bestAsk.Size = book.Asks[0].Size
	bestBid.Price = book.Bids[0].Price
	bestBid.Size = book.Bids[0].Size

	updateResp := UpdateResp{}
	var msg = []byte{}
	// listen for the incremental updates
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Fatal("Dydx c.ReadMessage() err | ", err)
			// TO DO: change from log.Fatal and implement reconnect
			// c, err = attemptReconnect(initResp)
			// subscribeChannel(c, initResp, payload)
			// if err != nil {
			// 	log.Fatal("Dydx c.ReadMessage() err | ", err)
			// }
		} else if !bytes.Equal(nil, msg) {
			err := json.Unmarshal(msg, &updateResp)
			if err != nil {
				log.Fatal("json.Unmarshal() error | ", err)
			} else {
				wsReceived := time.Now()
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bidVolume decimal.Decimal
				var askVolume decimal.Decimal
				var err error

				// get new asks and bids from websocket update
				var newBids []Bid
				for _, el := range updateResp.Contents.Bids {
					var bid Bid
					bid.Price, err = decimal.NewFromString(el.([]interface{})[0].(string))
					if err != nil {
						log.Fatal("decimal.NewFromString error | ", err)
					}
					bid.Size, err = decimal.NewFromString(el.([]interface{})[1].(string))
					if err != nil {
						log.Fatal("decimal.NewFromString error | ", err)
					}
					newBids = append(newBids, bid)
				}
				var newAsks []Ask
				for _, el := range updateResp.Contents.Asks {
					var ask Ask
					ask.Price, err = decimal.NewFromString(el.([]interface{})[0].(string))
					if err != nil {
						log.Fatal("decimal.NewFromString error | ", err)
					}
					ask.Size, err = decimal.NewFromString(el.([]interface{})[1].(string))
					if err != nil {
						log.Fatal("decimal.NewFromString error | ", err)
					}
					newAsks = append(newAsks, ask)
				}

				// update book with new bids and asks
				if len(newBids) > 0 {
					for _, newBid := range newBids {
						for i, bid := range book.Bids {
							//not sure why i needed to add the !size.equal to zero here but best bid/ask
							//were updating with zero volume levels, may mean this logic is broken
							if newBid.Price.GreaterThan(bid.Price) && !newBid.Size.Equal(decimal.Zero) {
								book.Bids = append([]Bid{newBid}, book.Bids...)
								break
							} else if newBid.Price.Equal(bid.Price) {
								if newBid.Size.Equal(decimal.Zero) { // remove from book where new size is zero
									book.Bids = append(book.Bids[:i], book.Bids[i+1:]...)
									break
								} else { // update with new size
									book.Bids[i].Size = newBid.Size
									break
								}
							}
						}
					}
				}
				if len(newAsks) > 0 {
					for _, newAsk := range newAsks {
						for i, ask := range book.Asks {
							//not sure why i needed to add the !size.equal to zero here but best bid/ask
							//were updating with zero volume levels, may mean this logic is broken
							if newAsk.Price.LessThan(ask.Price) && !newAsk.Size.Equal(decimal.Zero) {
								book.Asks = append([]Ask{newAsk}, book.Asks...)
								break
							} else if newAsk.Price.Equal(ask.Price) {
								if newAsk.Size.Equal(decimal.Zero) { // remove from book where new size is zero
									book.Asks = append(book.Asks[:i], book.Asks[i+1:]...)
									break
								} else { // update with new size
									book.Asks[i].Size = newAsk.Size
									break
								}
							}
						}
					}
				}

				if !book.Bids[0].Equal(bestBid) {
					bestBid = book.Bids[0]
					askVolume = bestAsk.Size
					ask = bestAsk.Price
					bid = bestBid.Price
					bidVolume = bestBid.Size
					updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
				}
				if !book.Asks[0].Equal(bestAsk) {
					bestAsk = book.Asks[0]
					askVolume = bestAsk.Size
					ask = bestAsk.Price
					bid = bestBid.Price
					bidVolume = bestBid.Size
					updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
				}
			}
		}
	}
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.dydx.exchange/v3/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
		// TO DO: implement redial and change to log.Println above
		// c = attemptRedial()
	}
	err = c.ReadJSON(&initResp)
	if err != nil {
		log.Fatal(err)
	} else if initResp["type"] != "connected" {
		log.Fatal("subscribe error | ", initResp)
	}
	log.Println("Connected.")
	return c, err
}

func subscribeChannel(c *websocket.Conn, initResp map[string]interface{}, payload string) []byte {
	var err error
	var msg = []byte{}
	c.WriteMessage(1, []byte(payload))
	_, msg, err = c.ReadMessage()
	if err != nil {
		log.Fatal("Dydx c.ReadMessage() err | ", err)
	} else {
		err := json.Unmarshal(msg, &initResp)
		if err != nil {
			log.Fatal(err)
		} else if initResp["type"] == "error" {
			log.Fatal("subscribe error | ", initResp)
		} else {
			log.Println("Subscribed.")
		}
	}
	return msg
}
