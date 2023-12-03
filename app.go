package main

import (
	"context"
	"fmt"
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

func (a *App) ConnectDydxWebsocket(pair string) {
	channelDydx := make(chan data.Spread)
	go dydx.GetSpread(channelDydx, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelDydx:
			runtime.EventsEmit(a.ctx, "spreadData", "DYDX", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectHyperliquidxWebsocket(pair string) {
	channelHyperliquidx := make(chan data.Spread)
	go hyperliquidx.GetSpread(channelHyperliquidx, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelHyperliquidx:
			runtime.EventsEmit(a.ctx, "spreadData", "HyperliquidX", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
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

func (a *App) ConnectKrakenFuturesWebsocket(pair string) {
	channelKrakenFutures := make(chan data.Spread)
	go krakenfutures.GetSpread(channelKrakenFutures, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelKrakenFutures:
			runtime.EventsEmit(a.ctx, "spreadData", "Kraken (Futures)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectOkxSpotWebsocket(pair string) {
	channelOkxSpot := make(chan data.Spread)
	go okxspot.GetSpread(channelOkxSpot, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelOkxSpot:
			runtime.EventsEmit(a.ctx, "spreadData", "Okx", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectOkxSwapsWebsocket(pair string) {
	channelOkxSwaps := make(chan data.Spread)
	go okxswaps.GetSpread(channelOkxSwaps, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelOkxSwaps:
			runtime.EventsEmit(a.ctx, "spreadData", "Okx (Swaps)", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
	}
}

func (a *App) ConnectUpbitWebsocket(pair string) {
	channelUpbit := make(chan data.Spread)
	go upbit.GetSpread(channelUpbit, pair)
	var spreadData = data.Spread{}
	for {
		select {
		case spreadData = <-channelUpbit:
			runtime.EventsEmit(a.ctx, "spreadData", "Upbit", spreadData.BidVolume, spreadData.Bid, spreadData.Ask, spreadData.AskVolume)
		}
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
	switch s, err := fetchFuncMap[exchange](); err {
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
}
