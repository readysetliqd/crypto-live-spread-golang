package bitgetfutures

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

// Accepts a channel and a pair. Connects to websocket api to updated channel
// with the live bid/ask spread for that pair
func GetSpread(updateChannel chan data.Spread, pair string) {
	// Bitget futures REST API returns pairs with suffixes eg _UMCBL, however
	// these pairs are only acceptable streams to subscribe in private streams
	// The public websocket streams accept base pairs eg BTCUSDT_UMCBL -> BTCUSDT
	pair = strings.Split(pair, "_")[0]
	defer close(updateChannel)
	var initResp map[string]interface{}
	var err error
	c, err := dial(initResp)
	if err != nil {
		log.Fatal("Failed initial connect | ", err)
	}
	defer c.Close()

	payload := "{\"op\": \"subscribe\", \"args\": [{\"instType\": \"mc\",\"channel\": \"ticker\",\"instId\": \"" + pair + "\"}]}"
	log.Println(payload)
	subscribeChannel(c, initResp, payload)

	// Create a ticker that fires every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	resp := map[string]interface{}{}
	// listen for the incremental updates
	for {
		select {
		case <-ticker.C:
			sendPing(c)
		default:
			_, msg, err := c.ReadMessage()
			msgString := string(msg[:])
			if err != nil {
				log.Fatal("Bitget c.ReadMessage() err | ", err)
				// TO DO: implement reconnect
				// c, err = attemptReconnect(initResp)
			} else if !bytes.Equal(nil, msg) && msgString != "pong" { //not a pong message
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

					bid, err = decimal.NewFromString(respInfc.(map[string]interface{})["bestBid"].(string))
					if err != nil {
						log.Fatal(err)
					}
					ask, err = decimal.NewFromString(respInfc.(map[string]interface{})["bestAsk"].(string))
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
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.bitget.com/mix/v1/stream", nil)
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
