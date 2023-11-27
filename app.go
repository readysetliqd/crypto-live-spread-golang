package main

import (
	"context"
	"fmt"
	"log"
	"slices"

	binancecoinm "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-coinm-futures"
	binancespot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-spot"
	binanceusdm "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binance-usdm-futures"
	binanceusspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/binanceus-spot"
	bitgetfutures "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bitget-futures"
	bitgetspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bitget-spot"
	bybitfutures "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bybit-futures"
	bybitspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/bybit-spot"
	coinbasespot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/coinbase-spot"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Selected market: %s", name)
}

func (a *App) ConnectKrakenSpotWebsocket(pair string) {
	var quit chan struct{} //nil
	channelKrakenSpot := make(chan krakenspot.Spread)
	go krakenspot.KrakenSpread(channelKrakenSpot, pair)
	var spreadData = krakenspot.Spread{}
	for {
		select {
		case spreadData = <-channelKrakenSpot:
			log.Println(spreadData)
			runtime.EventsEmit(a.ctx, "spreadData", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		case <-quit:
			return
		}
	}
}

func (a *App) FetchBinanceSpotPairs() []string {
	s := binancespot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBinanceUsdmPairs() []string {
	s := binanceusdm.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBinanceCoinmPairs() []string {
	s := binancecoinm.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBinanceUsSpotPairs() []string {
	s := binanceusspot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBitgetSpotPairs() []string {
	s := bitgetspot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBitgetFuturesPairs() []string {
	s := bitgetfutures.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBybitSpotPairs() []string {
	s := bybitspot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchBybitFuturesPairs() []string {
	s := bybitfutures.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchCoinbaseSpotPairs() []string {
	s := coinbasespot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchKrakenSpotPairs() []string {
	s := krakenspot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchKrakenFuturesPairs() []string {
	s := krakenfutures.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchOkxSpotPairs() []string {
	s := okxspot.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchOkxSwapsPairs() []string {
	s := okxswaps.FetchPairs()
	slices.Sort(s)
	return s
}

func (a *App) FetchUpbitSpotPairs() []string {
	s := upbit.FetchPairs()
	slices.Sort(s)
	return s
}
