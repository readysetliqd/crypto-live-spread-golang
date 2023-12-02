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

// Roughly translated from the api docs at https://docs.upbit.com/v1.4.0/reference/connection#pingpong
//
// By default, the server terminates the WebSocket Connection with Idle Timeout
// after about 120 seconds without data being displayed or sent.
// To prevent this, you can send a PING message from the client to the server
// to maintain the connection and view the state of the WebSocket and the
// WebSocket connection status.
//
// Currently, the Bit OpenAPI WebSocket server is in progress of preparing a
// PING Frame response, and the status of the server can be viewed through the
// PING request/PONG response (response frame to PING) as a property of the
// client. For configuration information, please refer to the relevant client
// development document. (For most libraries, ping is likely built-in.)
const timerDelay = 119 * time.Second

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

	// Create a ticker that fires every const timerDelay seconds
	timeTicker := time.NewTicker(timerDelay)
	defer timeTicker.Stop()

	resp := map[string]interface{}{}
	// listen for incremental updates and timeTicker countdown
	for {
		select {
		case <-timeTicker.C:
			sendPing(c)
		default:
			timeTicker.Reset(timerDelay)
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Okx c.ReadMessage() err | ", err)
				// TO DO: change from log.Fatal and implement reconnect
				// c, err = attemptReconnect(initResp)
			} else if !bytes.Equal(nil, msg) && msgType != websocket.PongMessage {
				err := json.Unmarshal(msg, &resp)
				if err != nil {
					log.Fatal(err)
				} else {
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
}

func sendPing(c *websocket.Conn) {
	payload := ""
	c.WriteMessage(websocket.PingMessage, []byte(payload))
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
