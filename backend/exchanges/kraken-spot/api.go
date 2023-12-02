package krakenspot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() []string {
	res, err := http.Get("https://api.kraken.com/0/public/AssetPairs")
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
	for _, el := range resp["result"].(map[string]interface{}) {
		//KrakenPairsList = append(KrakenPairsList, k)
		//KrakenPairsAltList = append(KrakenPairsAltList, el.(map[string]interface{})["wsname"].(string))
		pairs = append(pairs, el.(map[string]interface{})["wsname"].(string))
	}
	// DEBUG
	// sort.Slice(KrakenPairsAltList, func(i, j int) bool {
	// 	return KrakenPairsAltList[i] < KrakenPairsAltList[j]
	// })
	// log.Println(KrakenPairsList)
	return pairs
}
