package binancespot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	"github.com/shopspring/decimal"
)

// Accepts a channel and a pair. Connects to websocket api to updated channel
// with the live bid/ask spread for that pair
func GetSpread(updateChannel chan data.Spread, pair string) {
	defer close(updateChannel)
	var init_resp map[string]interface{}
	var err error
	pair = strings.ToLower(pair)
	c, err := dialBinance(init_resp, pair)
	if err != nil {
		log.Fatal("Failed initial connect | ", err)
	}
	defer c.Close()
	// payload := "{\"method\": \"SUBSCRIBE\", \"params\":[\"" + pair + "@bookTicker\"],\"id\": 1}"

	// log.Println(payload)
	// subscribeChannel(c, init_resp, payload)

	var resp interface{}
	var msg = []byte{}
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Println("Binance c.ReadMessage() err | ", err)
			// implement reconnection attempt
			// c, err = attemptReconnectBinance(init_resp)
			// subscribeChannel(c, init_resp, payload)
			// if err != nil {
			//	log.Fatal("Binance c.ReadMessage() err | ", err)
			//}
		} else if !bytes.Equal(nil, msg) { //not a ping frame message
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal("json.Unmarshal error | ", err)
			} else {
				wsReceived := time.Now()
				var bidVolume decimal.Decimal
				var bid decimal.Decimal
				var ask decimal.Decimal
				var askVolume decimal.Decimal
				respPair := resp.(map[string]interface{})["s"].(string)
				respPair = strings.ToLower(respPair)
				if pair != respPair {
					log.Fatal("Subscribed pair doesn't match input | ", pair, respPair)
				}
				bidVolume, err = decimal.NewFromString(resp.(map[string]interface{})["B"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				bid, err = decimal.NewFromString(resp.(map[string]interface{})["b"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				ask, err = decimal.NewFromString(resp.(map[string]interface{})["a"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				askVolume, err = decimal.NewFromString(resp.(map[string]interface{})["A"].(string))
				if err != nil {
					log.Fatal("decimal.NewFromString error | ", err)
				}
				updateChannel <- data.Spread{WsReceived: wsReceived, BidVolume: bidVolume, Bid: bid, Ask: ask, AskVolume: askVolume}
			}
		}
	}
}

func dialBinance(init_resp map[string]interface{}, pair string) (*websocket.Conn, error) {
	url := "wss://data-stream.binance.vision/ws/" + pair + "@bookTicker"
	c, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Println("dial error | ", err)
		// implement redial attempts
		// c = attemptRedialBinance()
	}

	err = c.ReadJSON(&init_resp)
	if err != nil {
		log.Fatal("ReadJSON() error | ", err)
	}
	log.Println("Connected.")
	return c, err
}
