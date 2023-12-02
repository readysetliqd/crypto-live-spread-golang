package okxspot

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

// From the docs at https://www.okx.com/docs-v5/en/#overview-websocket-connect:
// To keep the connection stable:
//
// 1. Set a timer of N seconds whenever a response message is received, where N
// is less than 30.
//
// 2. If the timer is triggered, which means that no new message is received
// within N seconds, send the String 'ping'.
//
// 3. Expect a 'pong' as a response. If the response message is not received
// within N seconds, please raise an error or reconnect.
const timerDelay = 29 * time.Second

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

	payload := "{\"op\": \"subscribe\", \"args\": [{\"channel\": \"tickers\",\"instId\": \"" + pair + "\"}]}"
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
			msgString := string(msg[:])
			if err != nil {
				log.Fatal("Okx c.ReadMessage() err | ", err)
				// TO DO: implement reconnect
				// c, err = attemptReconnect(initResp)
			} else if !bytes.Equal(nil, msg) && msgString != "pong" {
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
					respInfc := resp["data"].([]interface{})[0]

					bid, err = decimal.NewFromString(respInfc.(map[string]interface{})["bidPx"].(string))
					if err != nil {
						log.Fatal(err)
					}
					ask, err = decimal.NewFromString(respInfc.(map[string]interface{})["askPx"].(string))
					if err != nil {
						log.Fatal(err)
					}
					bidVolume, err = decimal.NewFromString(respInfc.(map[string]interface{})["bidSz"].(string))
					if err != nil {
						log.Fatal(err)
					}
					askVolume, err = decimal.NewFromString(respInfc.(map[string]interface{})["askSz"].(string))
					if err != nil {
						log.Fatal(err)
					}
					updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
				}
			}
		}
	}
}

func sendPing(c *websocket.Conn) {
	payload := `ping`
	c.WriteMessage(1, []byte(payload))
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.okx.com:8443/ws/v5/public", nil)
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
	if err != nil {
		log.Fatal(err)
	} else if initResp["event"] == "error" {
		log.Fatal("event error", initResp)
	}
	log.Println(initResp)
}
