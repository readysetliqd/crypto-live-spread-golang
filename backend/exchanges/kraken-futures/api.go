package krakenfutures

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://futures.kraken.com/derivatives/api/v3/instruments/status")
	if err != nil {
		log.Fatal("http.Get error | ", err)
	}
	defer res.Body.Close()

	resp := map[string]interface{}{}
	var msg = []byte{}
	msg, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("io.ReadAll error | ", err)
	}
	var pairs []string
	json.Unmarshal(msg, &resp)
	for _, symbol := range resp["instrumentStatus"].([]interface{}) {
		pairs = append(pairs, symbol.(map[string]interface{})["tradeable"].(string))
	}
	return pairs
}