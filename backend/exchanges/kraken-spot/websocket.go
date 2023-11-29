package krakenspot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const maxInstReconnectAttempts = 5
const reconnectDelay = 30 // seconds
const redialDelay = 30    //seconds

var heartbeat = []byte{123, 34, 101, 118, 101, 110, 116, 34, 58, 34, 104, 101, 97, 114, 116, 98, 101, 97, 116, 34, 125}
var reconnect_attempts int = 0
var redial_attempts int = 0

// Accepts a channel and a pair. Connects to websocket api to updated channel
// with the live bid/ask spread for that pair
func GetSpread(update_channel chan Spread, pair string) {
	defer close(update_channel)
	var init_resp map[string]interface{}
	var err error
	c, err := dialKraken(init_resp)
	if err != nil {
		log.Fatal("Failed initial connect | ", err)
	}
	defer c.Close()

	payload := "{\"event\": \"subscribe\", \"pair\": [\"" + pair + "\"], \"subscription\": {\"name\": \"spread\"}}"
	log.Println(payload)
	subscribeChannel(c, init_resp, payload)

	resp := []interface{}{}
	var msg = []byte{}
	// listen for the incremental updates
	for {
		_, msg, err = c.ReadMessage()
		if err != nil {
			log.Println("Kraken c.ReadMessage() err | ", err)
			c, err = attemptReconnectKraken(init_resp)
			subscribeChannel(c, init_resp, payload)
			if err != nil {
				log.Fatal("Kraken c.ReadMessage() err | ", err)
			}
		} else if !bytes.Equal(heartbeat, msg) { //not a heartbeat message
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Fatal(err)
			} else {
				wsReceived := time.Now()
				var time decimal.Decimal
				var bid decimal.Decimal
				var ask decimal.Decimal
				var bid_volume decimal.Decimal
				var ask_volume decimal.Decimal
				var err error
				resp_interface := resp[1].([]interface{})

				bid, err = decimal.NewFromString(resp_interface[0].(string))
				if err != nil {
					log.Fatal(err)
				}
				ask, err = decimal.NewFromString(resp_interface[1].(string))
				if err != nil {
					log.Fatal(err)
				}
				time, err = decimal.NewFromString(resp_interface[2].(string))
				if err != nil {
					log.Fatal(err)
				}
				bid_volume, err = decimal.NewFromString(resp_interface[3].(string))
				if err != nil {
					log.Fatal(err)
				}
				ask_volume, err = decimal.NewFromString(resp_interface[4].(string))
				if err != nil {
					log.Fatal(err)
				}
				update_channel <- Spread{WsReceived: wsReceived, Time: time, Bid: bid, Ask: ask, BidVolume: bid_volume, AskVolume: ask_volume}
			}
		}
	}
}

func dialKraken(init_resp map[string]interface{}) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.kraken.com", http.Header{})
	if err != nil {
		log.Println("dial:", err)
		c = attemptRedialKraken()
	}

	err = c.ReadJSON(&init_resp)
	if err != nil {
		log.Fatal(err)
	}
	if !(init_resp["event"] == "systemStatus" && init_resp["status"] == "online") {
		log.Fatal(init_resp)
	} else {
		log.Println("Connected.")
	}
	return c, err
}

func attemptRedialKraken() *websocket.Conn {
	log.Println("Attempting redial...")
	if redial_attempts < maxInstReconnectAttempts {
		c, _, err := websocket.DefaultDialer.Dial("wss://ws.kraken.com", http.Header{})
		if err == nil {
			return c
		}
		redial_attempts++
		time.Sleep(1 * time.Second)
	} else {
		c, _, err := websocket.DefaultDialer.Dial("wss://ws.kraken.com", http.Header{})
		if err == nil {
			return c
		}
		redial_attempts++
		log.Printf("Too many redials, trying again in %d seconds", redialDelay)
		time.Sleep(redialDelay * time.Second)
	}
	return attemptRedialKraken()
}

func subscribeChannel(c *websocket.Conn, init_resp map[string]interface{}, payload string) {
	var err error
	c.WriteMessage(1, []byte(payload))
	err = c.ReadJSON(&init_resp)
	// TO DO:
	// add check that all pairs are subscribed, previously:
	// && init_resp["pair"] == "XBT/USD"
	// doesn't work with multiple pairs
	if err != nil {
		log.Fatal(err)
	} else if !(init_resp["event"] == "subscriptionStatus" && init_resp["status"] == "subscribed") {
		log.Fatal(init_resp)
	}
	log.Println(init_resp)
}

func attemptReconnectKraken(init_resp map[string]interface{}) (*websocket.Conn, error) {
	log.Println("Attempting reconnect...")
	var c *websocket.Conn
	var err error
	if reconnect_attempts < maxInstReconnectAttempts {
		c, err = dialKraken(init_resp)
		reconnect_attempts++
		if err != nil {
			log.Println("Reconnect err | ", err)
			time.Sleep(1 * time.Second)
			return attemptReconnectKraken(init_resp)
		}
	} else {
		c, err = dialKraken(init_resp)
		reconnect_attempts++
		if err != nil {
			log.Println("Reconnect err | ", err)
			log.Printf("Too many reconnects, trying again in %d seconds", reconnectDelay)
			time.Sleep(reconnectDelay * time.Second)
			return attemptReconnectKraken(init_resp)
		}
	}
	log.Println("Reconnect successful")
	reconnect_attempts = 0
	return c, err
}
