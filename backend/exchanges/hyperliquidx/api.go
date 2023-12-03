package hyperliquidx

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() ([]string, int) {
	posturl := "https://api.hyperliquid.xyz/info"

	values := map[string]string{"type": "meta"}
	json_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post(posturl, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
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
		for _, symbol := range resp["universe"].([]interface{}) {
			pairs = append(pairs, symbol.(map[string]interface{})["name"].(string))
		}
		return pairs, 0
	} else if res.StatusCode == http.StatusForbidden {
		return []string{}, http.StatusForbidden
	} else {
		return []string{}, 1
	}
}
