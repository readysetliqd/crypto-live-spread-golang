package main

import (
	"context"
	"fmt"
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
	channelBinanceSpot := make(chan data.Spread)
	go binancespot.GetSpread(channelBinanceSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBinanceUsdmWebsocket(pair string) {
	channelBinanceUsdm := make(chan data.Spread)
	go binanceusdm.GetSpread(channelBinanceUsdm, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceUsdm:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance (USD-M)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBinanceCoinmWebsocket(pair string) {
	channelBinanceCoinm := make(chan data.Spread)
	go binancecoinm.GetSpread(channelBinanceCoinm, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceCoinm:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance (COIN-M)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBinanceUsWebsocket(pair string) {
	channelBinanceUs := make(chan data.Spread)
	go binanceus.GetSpread(channelBinanceUs, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBinanceUs:
			runtime.EventsEmit(a.ctx, "spreadData", "Binance US", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBitgetSpotWebsocket(pair string) {
	channelBitgetSpot := make(chan data.Spread)
	go bitgetspot.GetSpread(channelBitgetSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBitgetSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Bitget", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBitgetFuturesWebsocket(pair string) {
	channelBitgetFutures := make(chan data.Spread)
	go bitgetfutures.GetSpread(channelBitgetFutures, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBitgetFutures:
			runtime.EventsEmit(a.ctx, "spreadData", "Bitget (Futures)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBybitSpotWebsocket(pair string) {
	channelBybitSpot := make(chan data.Spread)
	go bybitspot.GetSpread(channelBybitSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBybitSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Bybit", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectBybitFuturesWebsocket(pair string) {
	channelBybitFutures := make(chan data.Spread)
	go bybitfutures.GetSpread(channelBybitFutures, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelBybitFutures:
			runtime.EventsEmit(a.ctx, "spreadData", "Bybit (Futures)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectCoinbaseSpotWebsocket(pair string) {
	channelCoinbaseSpot := make(chan data.Spread)
	go coinbasespot.GetSpread(channelCoinbaseSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelCoinbaseSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Coinbase", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectKrakenSpotWebsocket(pair string) {
	channelKrakenSpot := make(chan data.Spread)
	go krakenspot.GetSpread(channelKrakenSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelKrakenSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Kraken", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
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
	s := binanceus.FetchPairs()
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
