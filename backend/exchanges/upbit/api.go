package upbit

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://api.upbit.com/v1/market/all")
	if err != nil {
		log.Fatal("http.Get error | ", err)
	}
	defer res.Body.Close()

	resp := []interface{}{}
	var msg = []byte{}
	msg, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("io.ReadAll error | ", err)
	}
	var pairs []string
	json.Unmarshal(msg, &resp)
	for _, symbol := range resp {
		pairs = append(pairs, symbol.(map[string]interface{})["market"].(string))
	}
	return pairs
}
