package binanceusdm

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://fapi.binance.com/fapi/v1/exchangeInfo")
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
	for _, symbol := range resp["symbols"].([]interface{}) {
		pairs = append(pairs, symbol.(map[string]interface{})["symbol"].(string))
	}
	return pairs
}
