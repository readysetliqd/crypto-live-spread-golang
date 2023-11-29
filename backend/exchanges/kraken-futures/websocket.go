package krakenfutures

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

var heartbeat = []byte{0} // TO DO: figure out what heartbeat frames are

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

	payload := "{\"event\": \"subscribe\",\"feed\": \"ticker\",\"product_ids\": [\"" + pair + "\"]}"
	log.Println(payload)
	subscribeChannel(c, initResp, payload)

	resp := map[string]interface{}{}
	var msg = []byte{}
	// listen for the incremental updates
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Fatal("Kraken c.ReadMessage() err | ", err)
			// TO DO: change from log.Fatal and implement reconnect
			// c, err = attemptReconnect(initResp)
			// subscribeChannel(c, initResp, payload)
			// if err != nil {
			// 	log.Fatal("Kraken c.ReadMessage() err | ", err)
			// }
		} else if !bytes.Equal(heartbeat, msg) { //not a heartbeat message
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal(err)
			} else {
				wsReceived := time.Now()
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bid_volume decimal.Decimal
				var ask_volume decimal.Decimal

				bid = decimal.NewFromFloat(resp["bid"].(float64))
				ask = decimal.NewFromFloat(resp["ask"].(float64))
				bid_volume = decimal.NewFromFloat(resp["bid_size"].(float64))
				ask_volume = decimal.NewFromFloat(resp["ask_size"].(float64))
				updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bid_volume, AskVolume: ask_volume}
			}
		}
	}
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://futures.kraken.com/ws/v1", nil)
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
	} else if initResp["event"] != "info" {
		log.Fatal("event error", initResp)
	}
	log.Println(initResp)
	err = c.ReadJSON(&initResp)
	// TO DO:
	// add check that all pairs are subscribed
	if err != nil {
		log.Fatal(err)
	} else if initResp["event"] != "subscribed" {
		log.Fatal("event error", initResp)
	}
	log.Println(initResp)
}
