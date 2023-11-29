package main

import (
	"context"
	"fmt"
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

func (a *App) ConnectBinanceSpotWebsocket(pair string) {
	var quit chan struct{} //nil
	channelBinanceSpot := make(chan binancespot.Spread)
	go binancespot.GetSpread(channelBinanceSpot, pair)
	var spreadData = binancespot.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		case <-quit:
			return
		}
	}
}

func (a *App) ConnectBinanceUsdmWebsocket(pair string) {
	var quit chan struct{} //nil
	channelBinanceUsdm := make(chan binanceusdm.Spread)
	go binanceusdm.GetSpread(channelBinanceUsdm, pair)
	var spreadData = binanceusdm.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceUsdm:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance (USD-M)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		case <-quit:
			return
		}
	}
}

func (a *App) ConnectBinanceCoinmWebsocket(pair string) {
	var quit chan struct{} //nil
	channelBinanceCoinm := make(chan binancecoinm.Spread)
	go binancecoinm.GetSpread(channelBinanceCoinm, pair)
	var spreadData = binancecoinm.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceCoinm:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance (COIN-M)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		case <-quit:
			return
		}
	}
}

func (a *App) ConnectKrakenSpotWebsocket(pair string) {
	var quit chan struct{} //nil
	channelKrakenSpot := make(chan krakenspot.Spread)
	go krakenspot.GetSpread(channelKrakenSpot, pair)
	var spreadData = krakenspot.Spread{}
	for {
		select {
		case spreadData = <-channelKrakenSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Kraken", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
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
