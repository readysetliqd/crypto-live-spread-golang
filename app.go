package main

import (
	"context"
	"fmt"
	"log"

	krakenspot "github.com/readysetliqd/crypto-live-spread-golang/backend/exchanges/kraken-spot"
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

func (a *App) FetchKrakenSpotPairs() map[string]string {
	return krakenspot.FetchPairs()
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
