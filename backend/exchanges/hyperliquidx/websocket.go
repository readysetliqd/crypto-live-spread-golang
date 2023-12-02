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

// From the docs at https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket#timeouts-and-heartbeats
//
// # Timeouts and Heartbeats
//
// The server will close any connection if it hasn't sent a message to it in
// the last 60 seconds. If you are subscribing to a channel that doesn't receive
// messages every 60 seconds, you can send heartbeat messages to keep your
// connection alive.
const timerDelay = 59 * time.Second

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

	// Create a ticker that fires every const timerDelay seconds
	timeTicker := time.NewTicker(timerDelay)
	defer timeTicker.Stop()

	resp := map[string]interface{}{}
	var msg = []byte{}
	// listen for incremental updates and timeTicker countdown
	for {
		select {
		case <-timeTicker.C:
			sendPing(c)
		default:
			timeTicker.Reset(timerDelay)
			_, msg, err = c.ReadMessage()
			if err != nil {
				log.Fatal("HyperliquidX c.ReadMessage() err | ", err)
				// TO DO: implement reconnect
				// c, err = attemptReconnect(initResp)
			} else if !bytes.Equal(nil, msg) {
				err := json.Unmarshal(msg, &resp)
				if err != nil {
					log.Fatal(err)
				} else {
					if resp["channel"].(string) == "pong" {
						continue
					}
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
}

func sendPing(c *websocket.Conn) {
	payload := `{ "method": "ping" }`
	c.WriteMessage(1, []byte(payload))
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
