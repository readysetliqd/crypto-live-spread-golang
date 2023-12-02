package okxswaps

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://www.okx.com/api/v5/public/instruments?instType=SWAP")
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
	for _, symbol := range resp["data"].([]interface{}) {
		pairs = append(pairs, symbol.(map[string]interface{})["instId"].(string))
	}
	return pairs
}
