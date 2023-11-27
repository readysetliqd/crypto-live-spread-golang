package krakenspot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchPairs() map[string]string {
	res, err := http.Get("https://api.kraken.com/0/public/AssetPairs")
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

	pairs := map[string]string{}
	json.Unmarshal(msg, &resp)
	for k, el := range resp["result"].(map[string]interface{}) {
		//KrakenPairsList = append(KrakenPairsList, k)
		//KrakenPairsAltList = append(KrakenPairsAltList, el.(map[string]interface{})["wsname"].(string))
		pairs[k] = el.(map[string]interface{})["wsname"].(string)
	}
	// DEBUG
	// sort.Slice(KrakenPairsAltList, func(i, j int) bool {
	// 	return KrakenPairsAltList[i] < KrakenPairsAltList[j]
	// })
	// log.Println(KrakenPairsList)
	return pairs
}
