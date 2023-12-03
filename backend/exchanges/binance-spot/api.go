package binancespot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() ([]string, int) {
	res, err := http.Get("https://data-api.binance.vision/api/v3/exchangeInfo")
	if err != nil {
		log.Fatal("http.Get error | ", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
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
		return pairs, 0
	} else if res.StatusCode == http.StatusForbidden {
		return []string{}, http.StatusForbidden
	} else {
		return []string{}, 1
	}
}
