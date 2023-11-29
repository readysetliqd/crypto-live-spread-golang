package hyperliquidx

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

	payload := `
	{
		"method": "subscribe",
		"subscription": { "type": "l2Book", "coin": "` + pair + `" }
	}
	`
	log.Println(payload)
	subscribeChannel(c, initResp, payload)

	resp := map[string]interface{}{}
	var msg = []byte{}
	// listen for the incremental updates
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Fatal("Bitget c.ReadMessage() err | ", err)
			// TO DO: change from log.Fatal and implement reconnect
			// c, err = attemptReconnect(initResp)
			// subscribeChannel(c, initResp, payload)
			// if err != nil {
			// 	log.Fatal("Bitget c.ReadMessage() err | ", err)
			// }
		} else if !bytes.Equal(pong, msg) { //not a pong message
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal(err)
			} else {
				wsReceived := time.Now()
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bidVolume decimal.Decimal
				var askVolume decimal.Decimal
				var err error
				respInfc := resp["data"].(map[string]interface{})
				bidsList := respInfc["levels"].([]interface{})[0]
				asksList := respInfc["levels"].([]interface{})[1]

				bestBid := bidsList.([]interface{})[0].(map[string]interface{})
				bestAsk := asksList.([]interface{})[0].(map[string]interface{})
				bidVolume, err = decimal.NewFromString(bestBid["sz"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				bid, err = decimal.NewFromString(bestBid["px"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				ask, err = decimal.NewFromString(bestAsk["px"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				askVolume, err = decimal.NewFromString(bestAsk["sz"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
			}
		}
	}
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.hyperliquid.xyz/ws", nil)
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
	} else if initResp["channel"] == "error" {
		log.Fatal("subscribe error | ", initResp)
	}
	log.Println(initResp)
}
