package bybitspot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://api.bybit.com/v5/market/instruments-info?category=spot")
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
	for _, symbol := range resp["result"].(map[string]interface{})["list"].([]interface{}) {
		pairs = append(pairs, symbol.(map[string]interface{})["symbol"].(string))
	}
	return pairs
}
