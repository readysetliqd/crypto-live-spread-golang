package main

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/readysetliqd/crypto-live-spread-golang/backend/data"
	binancecoinm "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-coinm-futures"
	binancespot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-spot"
	binanceusdm "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-usdm-futures"
	binanceus "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binanceus-spot"
	bitgetfutures "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bitget-futures"
	bitgetspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bitget-spot"
	bybitfutures "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bybit-futures"
	bybitspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bybit-spot"
	coinbasespot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/coinbase-spot"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/dydx"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/hyperliquidx"
	krakenfutures "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/kraken-futures"
	krakenspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/kraken-spot"
	okxspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/okx-spot"
	okxswaps "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/okx-swaps"
	"github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/upbit"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ConnectWebsocket(exchange string, pair string) {
	connectFuncMap := map[string]func(chan data.Spread, string){
		"Binance":          binancespot.GetSpread,
		"Binance (USD-M)":  binanceusdm.GetSpread,
		"Binance (COIN-M)": binancecoinm.GetSpread,
		"Binance US":       binanceus.GetSpread,
		"Bitget":           bitgetspot.GetSpread,
		"Bitget (Futures)": bitgetfutures.GetSpread,
		"Bybit":            bybitspot.GetSpread,
		"Bybit (Futures)":  bybitfutures.GetSpread,
		"Coinbase":         coinbasespot.GetSpread,
		"DYDX":             dydx.GetSpread,
		"HyperliquidX":     hyperliquidx.GetSpread,
		"Kraken":           krakenspot.GetSpread,
		"Kraken (Futures)": krakenfutures.GetSpread,
		"Okx":              okxspot.GetSpread,
		"Okx (Swaps)":      okxswaps.GetSpread,
		"Upbit":            upbit.GetSpread,
	}
	if connectFunc, exists := connectFuncMap[exchange]; exists {
		dataChan := make(chan data.Spread)
		quit := make(chan bool)
		go connectFunc(dataChan, pair)
		go func() {
			var spreadData = data.Spread{}
			for {
				select {
				case spreadData = <-dataChan:
					runtime.EventsEmit(a.ctx, "spreadData", exchange, spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
				case <-quit:
					return
				}
			}
		}()
		runtime.EventsOn(a.ctx, "stopGoroutine", func(optionalData ...interface{}) {
			if len(optionalData) > 0 {
				if e, ok := optionalData[0].(string); ok && e == exchange {
					log.Println("Quitting goroutine exchange: ", e)
					close(quit)
					quit = make(chan bool)
				}
			}
		})
	} else {
		log.Fatal("ConnectWebsocket() error | Function not found for exchange: ", exchange)
	}
}

func (a *App) FetchPairs(exchange string) []string {
	fetchFuncMap := map[string]func() ([]string, int){
		"Binance":          binancespot.FetchPairs,
		"Binance (USD-M)":  binanceusdm.FetchPairs,
		"Binance (COIN-M)": binancecoinm.FetchPairs,
		"Binance US":       binanceus.FetchPairs,
		"Bitget":           bitgetspot.FetchPairs,
		"Bitget (Futures)": bitgetfutures.FetchPairs,
		"Bybit":            bybitspot.FetchPairs,
		"Bybit (Futures)":  bybitfutures.FetchPairs,
		"Coinbase":         coinbasespot.FetchPairs,
		"DYDX":             dydx.FetchPairs,
		"HyperliquidX":     hyperliquidx.FetchPairs,
		"Kraken":           krakenspot.FetchPairs,
		"Kraken (Futures)": krakenfutures.FetchPairs,
		"Okx":              okxspot.FetchPairs,
		"Okx (Swaps)":      okxswaps.FetchPairs,
		"Upbit":            upbit.FetchPairs,
	}
	if fetchFunc, exists := fetchFuncMap[exchange]; exists {
		switch s, err := fetchFunc(); err {
		case 0:
			slices.Sort(s)
			return s
		case 1:
			runtime.EventsEmit(a.ctx, "Unspecified error")
			return s
		case http.StatusForbidden:
			runtime.EventsEmit(a.ctx, "HTTP Forbidden")
			return s
		default:
			return s
		}
	} else {
		log.Fatal("FetchPairs() error | Function not found for exchange: ", exchange)
		return []string{}
	}
}
