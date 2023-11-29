package upbit

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

var pong = []byte{0} // TO DO: figure out what pong frames are

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

	payload := "[{\"ticket\": \"1001\"},{\"type\": \"orderbook\",\"codes\": [\"" + pair + ".1\"]},{\"format\": \"DEFAULT\"}]"
	log.Println(payload)
	subscribeChannel(c, initResp, payload)

	resp := map[string]interface{}{}
	var msg = []byte{}
	// listen for the incremental updates
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Fatal("Okx c.ReadMessage() err | ", err)
			// TO DO: change from log.Fatal and implement reconnect
			// c, err = attemptReconnect(initResp)
			// subscribeChannel(c, initResp, payload)
			// if err != nil {
			// 	log.Fatal("Okx c.ReadMessage() err | ", err)
			// }
		} else if !bytes.Equal(pong, msg) { //not a pong message
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println(resp)
				wsReceived := time.Now()
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bidVolume decimal.Decimal
				var askVolume decimal.Decimal
				respInfc := resp["orderbook_units"].([]interface{})[0].(map[string]interface{})

				bid = decimal.NewFromFloat(respInfc["bid_price"].(float64))
				ask = decimal.NewFromFloat(respInfc["ask_price"].(float64))
				bidVolume = decimal.NewFromFloat(respInfc["bid_size"].(float64))
				askVolume = decimal.NewFromFloat(respInfc["ask_size"].(float64))
				updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
			}
		}
	}
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.upbit.com/websocket/v1", nil)
	if err != nil {
		log.Fatal("dial:", err)
		// TO DO: implement redial and change to log.Println above
		// c = attemptRedial()
	}
	log.Println("Connected.")
	return c, err
}

func subscribeChannel(c *websocket.Conn, initResp map[string]interface{}, payload string) {
	var err error
	c.WriteMessage(1, []byte(payload))
	err = c.ReadJSON(&initResp)
	// TO DO:
	// add check that all pairs are subscribed
	if err != nil {
		log.Fatal(err)
	} else if initResp["event"] == "error" {
		log.Fatal("event error", initResp)
	}
	log.Println(initResp)
}
