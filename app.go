package main

import (
	"context"
	"fmt"
	"log"

	krakenspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/kraken"
	"github.com/wailsapp/wails"
)

// App struct
type App struct {
	ctx     context.Context
	runtime *wails.Runtime
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context, runtime *wails.Runtime) {
	a.ctx = ctx
	a.runtime = runtime
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Selected market: %s", name)
}

func (a *App) FetchKrakenSpotPairs() map[string]string {
	return krakenspot.FetchPairs()
}

func (a *App) ConnectKrakenSpotWebsocket(pair string) {
	var quit chan struct{} //nil
	channel_kraken_spot := make(chan krakenspot.Spread)
	go krakenspot.KrakenSpread(channel_kraken_spot, pair)
	var spread_data = krakenspot.Spread{}
	for {
		select {
		case spread_data = <-channel_kraken_spot:
			log.Println(spread_data)
		case <-quit:
			return
		}
	}
}
