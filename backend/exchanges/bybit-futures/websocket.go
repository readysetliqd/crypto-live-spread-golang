package bybitfutures

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

	payload := "{\"op\": \"subscribe\", \"args\": [\"orderbook.1." + pair + "\"]}"
	log.Println(payload)
	subscribeChannel(c, initResp, payload)

	// Create a ticker that fires every 20 seconds
	timeTicker := time.NewTicker(20 * time.Second)
	defer timeTicker.Stop()

	var msg = []byte{}
	// listen for the incremental updates
	var prevBid decimal.Decimal
	var prevAsk decimal.Decimal
	var prevBidVolume decimal.Decimal
	var prevAskVolume decimal.Decimal
	for {
		select {
		case <-timeTicker.C:
			sendPing(c)
		default:
			resp := map[string]interface{}{}
			_, msg, err = c.ReadMessage()
			if err != nil {
				log.Fatal("Bybit c.ReadMessage() err | ", err)
				// TO DO: implement reconnect
				// c, err = attemptReconnect(initResp)
			}
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal(err)
			}
			if resp["op"] != nil {
				if resp["op"].(string) == "ping" { // pong message
					continue
				}
			} else if !bytes.Equal(msg, nil) {
				wsReceived := time.Now()
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bidVolume decimal.Decimal
				var askVolume decimal.Decimal
				var err error
				respInfc := resp["data"].(map[string]interface{})
				if resp["type"].(string) == "snapshot" {
					bid, err = decimal.NewFromString(respInfc["b"].([]interface{})[0].([]interface{})[0].(string))
					if err != nil {
						log.Fatal(err)
					}
					ask, err = decimal.NewFromString(respInfc["a"].([]interface{})[0].([]interface{})[0].(string))
					if err != nil {
						log.Fatal(err)
					}
					bidVolume, err = decimal.NewFromString(respInfc["b"].([]interface{})[0].([]interface{})[1].(string))
					if err != nil {
						log.Fatal(err)
					}
					askVolume, err = decimal.NewFromString(respInfc["a"].([]interface{})[0].([]interface{})[1].(string))
					if err != nil {
						log.Fatal(err)
					}
					prevBidVolume = bidVolume
					prevBid = bid
					prevAsk = ask
					prevAskVolume = askVolume
					updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}
				} else if resp["type"].(string) == "delta" {
					if len(respInfc["a"].([]interface{})) > 0 {
						for _, el := range respInfc["a"].([]interface{}) {
							if el.([]interface{})[1].(string) != "0" {
								ask, err = decimal.NewFromString(el.([]interface{})[0].(string))
								if err != nil {
									log.Fatal(err)
								}
								askVolume, err = decimal.NewFromString(el.([]interface{})[1].(string))
								if err != nil {
									log.Fatal(err)
								}
								prevAsk = ask
								prevAskVolume = askVolume
							}
						}
					} else {
						ask = prevAsk
						askVolume = prevAskVolume
					}
					if len(respInfc["b"].([]interface{})) > 0 {
						for _, el := range respInfc["b"].([]interface{}) {
							if el.([]interface{})[1].(string) != "0" {
								bid, err = decimal.NewFromString(el.([]interface{})[0].(string))
								if err != nil {
									log.Fatal(err)
								}
								bidVolume, err = decimal.NewFromString(el.([]interface{})[1].(string))
								if err != nil {
									log.Fatal(err)
								}
								prevBid = bid
								prevBidVolume = bidVolume
							}
						}
					} else {
						bid = prevBid
						bidVolume = prevBidVolume
					}
					updateChannel <- data.Spread{WsReceived: wsReceived, Bid: bid, Ask: ask, BidVolume: bidVolume, AskVolume: askVolume}

				}
			}
		}
	}
}

func sendPing(c *websocket.Conn) {
	payload := `{"op": "ping"}`
	c.WriteMessage(1, []byte(payload))
}

func dial(initResp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/public/linear", nil)
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
	}
	// not sure why bybit spot has an initial connection response but bybit futures doesn't
	// } else if !initResp["success"].(bool) {
	// 	log.Fatal("event error", initResp)
	// }
}
