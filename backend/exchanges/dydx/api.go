package dydx

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://api.dydx.exchange/v3/markets")
	if err != nil {
		log.Fatal("http.Get error | ", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return []string{"Error accessing endpoint, check if your IP is geoblocked"}
	}

	resp := map[string]interface{}{}
	var msg = []byte{}
	msg, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("io.ReadAll error | ", err)
	}
	var pairs []string
	json.Unmarshal(msg, &resp)
	for _, symbol := range resp["markets"].(map[string]interface{}) {
		pairs = append(pairs, symbol.(map[string]interface{})["market"].(string))
	}
	return pairs
}
