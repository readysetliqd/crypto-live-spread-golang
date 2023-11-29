<script lang="ts">
  import { EventsOn } from '../wailsjs/runtime'
  import { onMount } from 'svelte'
  import { Svroller } from "svrollbar"
  import * as go from '../wailsjs/go/main/App.js'
    import { group_outros } from 'svelte/internal';

  let coin: string = ''
  let denom: string = ''
  let pairsPlaceholder: string[] = []
  
  onMount(() => {
    EventsOn("spreadData", (exchange, bidVolume, bid, ask, askVolume) => {
      exchanges[exchange].bidVolume = bidVolume
      exchanges[exchange].bid = bid
      exchanges[exchange].ask = ask
      exchanges[exchange].askVolume = askVolume
    })
  })

  $: {coin, denom, 
    Object.keys(exchanges).forEach((exchange) => {
      exchanges[exchange].recommended = getRecommendedPairs(exchange);
    });
  }

  function getRecommendedPairs(exchange) {
    let coinInput: string
    let denomInput: string
    if ((exchange == "Kraken" || exchange == "Kraken (Futures)") && (coin == "BTC" || coin == "")) {
      coinInput = "XBT"
    } else if (coin == "") { 
      coinInput = "BTC"
    } else {
      coinInput = coin
    }
    if ((exchange == "Kraken" || exchange == "Kraken (Futures)") && coin == "BTC") {
      denomInput == "XBT"
    } else if (denom == "") {
      denomInput = "USD"
    } else {
      denomInput = denom
    }
    return exchanges[exchange].pairs.filter(pair => pair.includes(coinInput) && pair.includes(denomInput))
  }

  function fetchPairs(exchange) {
    console.log('fetchPairs called', exchange)
    switch (exchange) {
      case 'Binance':
        fetchBinanceSpotPairs(exchange)
        break
      case 'Binance (USD-M)':
        fetchBinanceUsdmPairs(exchange)
        break
      case 'Binance (COIN-M)':
        fetchBinanceCoinmPairs(exchange)
        break
      case 'Binance US':
        fetchBinanceUsSpotPairs(exchange)
        break
      case 'Bitget':
        fetchBitgetSpotPairs(exchange)
        break
      case 'Bitget (Futures)':
        fetchBitgetFuturesPairs(exchange)
        break
      case 'Bybit':
        fetchBybitSpotPairs(exchange)
        break
      case 'Bybit (Futures)':
        fetchBybitFuturesPairs(exchange)
        break
      case 'Coinbase':
        fetchCoinbaseSpotPairs(exchange)
        break
      case 'Kraken':
        fetchKrakenSpotPairs(exchange)
        break
      case 'Kraken (Futures)':
        fetchKrakenFuturesPairs(exchange)
        break
      case 'Okx':
        fetchOkxSpotPairs(exchange)
        break
      case 'Okx (Swaps)':
        fetchOkxSwapsPairs(exchange)
        break
      case 'Upbit':
        fetchUpbitSpotPairs(exchange)
        break
    }
  }
  function connectWebsocket(exchange, pair) {
    console.log('connectWebsocket called', exchange)
    switch (exchange) {
      case 'Binance':
        connectBinanceSpotWebsocket(pair)
        break
      case 'Binance (USD-M)':
        connectBinanceUsdmWebsocket(pair)
        break
      case 'Binance (COIN-M)':
        connectBinanceCoinmWebsocket(pair)
        break
      case 'Binance US':
        connectBinanceUsWebsocket(pair)
        break
      case 'Bitget':
        connectBitgetSpotWebsocket(pair)
        break
      case 'Bitget (Futures)':
        connectBitgetFuturesWebsocket(pair)
        break
      case 'Bybit':
        connectBybitSpotWebsocket(pair)
        break
      case 'Bybit (Futures)':
        connectBybitFuturesWebsocket(pair)
        break
      case 'Coinbase':
        connectCoinbaseSpotWebsocket(pair)
        break
      case 'Kraken':
        connectKrakenSpotWebsocket(pair)
        break
      case 'Kraken (Futures)':
        connectKrakenFuturesWebsocket(pair)
        break
      case 'Okx':
        connectOkxSpotWebsocket(pair)
        break
      case 'Okx (Swaps)':
        connectOkxSwapsWebsocket(pair)
        break
      case 'Upbit':
        connectUpbitWebsocket(pair)
        break
    }
  }
  //#region fetchExchangePairs(exchange) functions
  async function fetchBinanceSpotPairs(exchange): Promise<void> {
    console.log('fetchBinanceSpotPairs called')
    const result = await go.FetchBinanceSpotPairs()
    console.log('FetchBinanceSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBinanceUsdmPairs(exchange): Promise<void> {
    console.log('fetchBinanceUsdmPairs called')
    const result = await go.FetchBinanceUsdmPairs()
    console.log('FetchBinanceUsdmPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBinanceCoinmPairs(exchange): Promise<void> {
    console.log('fetchBinanceCoinmPairs called')
    const result = await go.FetchBinanceCoinmPairs()
    console.log('FetchBinanceCoinmPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBinanceUsSpotPairs(exchange): Promise<void> {
    console.log('fetchBinanceUsSpotPairs called')
    const result = await go.FetchBinanceUsSpotPairs()
    console.log('FetchBinanceUsSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBitgetSpotPairs(exchange): Promise<void> {
    console.log('fetchBitgetSpotPairs called')
    const result = await go.FetchBitgetSpotPairs()
    console.log('FetchBitgetSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBitgetFuturesPairs(exchange): Promise<void> {
    console.log('fetchBitgetFuturesPairs called')
    const result = await go.FetchBitgetFuturesPairs()
    console.log('FetchBitgetFuturesPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBybitSpotPairs(exchange): Promise<void> {
    console.log('fetchBybitSpotPairs called')
    const result = await go.FetchBybitSpotPairs()
    console.log('FetchBybitSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchBybitFuturesPairs(exchange): Promise<void> {
    console.log('fetchBybitFuturesPairs called')
    const result = await go.FetchBybitFuturesPairs()
    console.log('FetchBybitFuturesPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchCoinbaseSpotPairs(exchange): Promise<void> {
    console.log('fetchCoinbaseSpotPairs called')
    const result = await go.FetchCoinbaseSpotPairs()
    console.log('FetchCoinbaseSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchKrakenSpotPairs(exchange): Promise<void> {
    console.log('fetchKrakenSpotPairs called')
    const result = await go.FetchKrakenSpotPairs()
    console.log('FetchKrakenSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchKrakenFuturesPairs(exchange): Promise<void> {
    console.log('fetchKrakenFuturesPairs called')
    const result = await go.FetchKrakenFuturesPairs()
    console.log('FetchKrakenFuturesPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchOkxSpotPairs(exchange): Promise<void> {
    console.log('fetchOkxSpotPairs called')
    const result = await go.FetchOkxSpotPairs()
    console.log('FetchOkxSpotPairs result', result)
    exchanges[exchange].pairs = result
  }

  async function fetchOkxSwapsPairs(exchange): Promise<void> {
    console.log('fetchOkxSwapsPairs called')
    const result = await go.FetchOkxSwapsPairs()
    console.log('FetchOkxSwapsPairs result', result)
    exchanges[exchange].pairs = result
  }
  
  async function fetchUpbitSpotPairs(exchange): Promise<void> {
    console.log('fetchUpbitSpotPairs called')
    const result = await go.FetchUpbitSpotPairs()
    console.log('FetchUpbitSpotPairs result', result)
    exchanges[exchange].pairs = result
  }
  //#endregion
  //#region connectExchangeWebsocket(pair) functions
  async function connectKrakenSpotWebsocket(pair): Promise<void> {
    console.log('connectKrakenSpotWebsocket called')
    const result = await go.ConnectKrakenSpotWebsocket(pair)
  }

  async function connectBinanceSpotWebsocket(pair): Promise<void> {
    console.log('connectBinanceSpotWebsocket called')
    const result = await go.ConnectBinanceSpotWebsocket(pair)
  }

  async function connectBinanceUsdmWebsocket(pair): Promise<void> {
    console.log('connectBinanceUsdmWebsocket called')
    const result = await go.ConnectBinanceUsdmWebsocket(pair)
  }

  async function connectBinanceCoinmWebsocket(pair): Promise<void> {
    console.log('connectBinanceCoinmWebsocket called')
    const result = await go.ConnectBinanceCoinmWebsocket(pair)
  }

  async function connectBinanceUsWebsocket(pair): Promise<void> {
    console.log('connectBinanceUsWebsocket called')
    const result = await go.ConnectBinanceUsWebsocket(pair)
  }

  async function connectBitgetSpotWebsocket(pair): Promise<void> {
    console.log('connectBitgetSpotWebsocket() called')
    const result = await go.ConnectBitgetSpotWebsocket(pair)
  }

  async function connectBitgetFuturesWebsocket(pair): Promise<void> {
    console.log('connectBitgetFuturesWebsocket() called')
    const result = await go.ConnectBitgetFuturesWebsocket(pair)
  }

  async function connectBybitSpotWebsocket(pair): Promise<void> {
    console.log('connectBybitSpotWebsocket() called')
    const result = await go.ConnectBybitSpotWebsocket(pair)
  }

  async function connectBybitFuturesWebsocket(pair): Promise<void> {
    console.log('connectBybitFuturesWebsocket() called')
    const result = await go.ConnectBybitFuturesWebsocket(pair)
  }

  async function connectCoinbaseSpotWebsocket(pair): Promise<void> {
    console.log('connectCoinbaseSpotWebsocket() caled')
    const result = await go.ConnectCoinbaseSpotWebsocket(pair)
  }

  async function connectKrakenFuturesWebsocket(pair): Promise<void> {
    console.log('connectKrakenFuturesWebsocket() called')
    const result = await go.ConnectKrakenFuturesWebsocket(pair)
  }

  async function connectOkxSpotWebsocket(pair): Promise<void> {
    console.log('connectOkxSpotWebsocket() called')
    const result = await go.ConnectOkxSpotWebsocket(pair)
  }

  async function connectOkxSwapsWebsocket(pair): Promise<void> {
    console.log('connectOkxSwapsWebsocket() called')
    const result = await go.ConnectOkxSwapsWebsocket(pair)
  }

  async function connectUpbitWebsocket(pair): Promise<void> {
    console.log('connectUpbitWebsocket() called')
    const result = await go.ConnectUpbitWebsocket(pair)
  }
  //#endregion

  const selectPair= (event) => {
    const [selectedPair, exchange] = event.target.value.split(',')
    connectWebsocket(exchange, selectedPair)
  }

  let exchanges = {
    "Binance": {category: "Spot", name: "Binance", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Binance US": {category: "Spot", name: "Binance US", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Bitget": {category: "Spot", name: "Bitget", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Coinbase": {category: "Spot", name: "Coinbase", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Okx": {category: "Spot", name: "Okx", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Kraken": {category: "Spot", name: "Kraken", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Bybit": {category: "Spot", name: "Bybit", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Upbit": {category: "Spot", name: "Upbit", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Kraken (Futures)": {category: "Futures", name: "Kraken (Futures)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Binance (USD-M)": {category: "Futures", name: "Binance (USD-M)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Binance (COIN-M)": {category: "Futures", name: "Binance (COIN-M)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Bybit (Futures)": {category: "Futures", name: "Bybit (Futures)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Okx (Swaps)": {category: "Futures", name: "Okx (Swaps)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "Bitget (Futures)": {category: "Futures", name: "Bitget (Futures)", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "DYDX": {category: "DEX", name: "DYDX", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "GMX": {category: "DEX", name: "GMX", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
    "HyperliquidX": {category: "DEX", name: "HyperliquidX", checked: false, pairs: pairsPlaceholder, recommended: pairsPlaceholder, bidVolume: 0, bid: 0, ask: 0, askVolume: 0},
  }
  let exchangeNames = Object.keys(exchanges)

  // Define a function to toggle the checked status of an exchange
  function toggleExchange(exchange) {
    console.log("exchange checked", exchanges)
    fetchPairs(exchange.name)
    exchange.checked = !exchange.checked;
  }
</script>

<main>
  <div class="body">
  <div class="title">Enter coin to track below</div>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={coin} class="input" id="coin" type="text" placeholder="BTC"/>
  </div>
  <div class="title">Enter denominator to track below (if empty, defaults to USD)</div>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={denom} class="input" id="denom" type="text" placeholder="USD"/>
  </div>
  <div class="column">
    <Svroller height="100%" width=12rem>
      <ul class="list">
        {#each exchangeNames as exchangeName, i (exchangeName)}
        {#if i === 0 || exchanges[exchangeNames[i]].category !== exchanges[exchangeNames[i-1]].category}
        <li class="category">{exchanges[exchangeName].category}</li>
        {/if}
        <li class="item" on:click={() => toggleExchange(exchanges[exchangeName])} on:keypress={() => toggleExchange(exchanges[exchangeName])}>
          <input
            class="tickbox"
            type="checkbox"
            bind:checked={exchanges[exchangeName].checked}
          />
          {exchangeName}
        </li>
        {/each}
      </ul>
    </Svroller>
  </div>
    {#each exchangeNames as exchange}
    {#if exchanges[exchange].checked}
    <table class="table">
      <thead>
        <tr>
          <th class="th">{exchange}</th>
          <th class="th">Bid Volume</th>
          <th class="th">Bid Price</th>
          <th class="th">Ask Price</th>
          <th class="th">Ask Volume</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td class="td1">
            <select on:change={selectPair} class="dropdown">
              <option selected disabled>Select pair...</option>
              <optgroup label="Recommended">
                {#if exchanges[exchange].recommended.length != 0}
                {#each exchanges[exchange].recommended as pair}
                <option value="{pair},{exchange}">{pair}</option>
                {/each}
                {:else}
                <option value="None" disabled>None</option>
                {/if}
              </optgroup>
              <optgroup label="All Pairs">
                {#each exchanges[exchange].pairs as pair}
                <option value="{pair},{exchange}">{pair}</option>
                {/each}
              </optgroup>
              </select>
          </td>
          <td class="td2">{exchanges[exchange].bidVolume}</td>
          <td class="td2">{exchanges[exchange].bid}</td>
          <td class="td2">{exchanges[exchange].ask}</td>
          <td class="td2">{exchanges[exchange].askVolume}</td>
        </tr>
      </tbody>
    </table>
    {/if}
    {/each}
  </div>
  
<!--
  <div>
    <select bind:value={selectedAsset} on:change={selectAsset}>
      {#each krakenSpotPairs as pair}
              <option value="{pair}">{pair}</option>
      {/each}
    </select>
  </div>
  <div class="selected-asset" id="selected-asset">{selectedAsset}</div>
  <div class="kraken-spread" id="kraken-spread">{krakenSpread}</div>-->
</main>

<style>
  .body {
    margin-left: 12rem;
  }
  .title {
    height: 20px;
    line-height: 20px;
    margin-top: 1.5rem;
    margin-bottom: 0.5rem;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
    margin-top: 1rem;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .column {
    position: fixed;
    left: 0;
    top: 0;
    width: 12rem;
    height: 100%;
    background-color: #3b3b3b;
    box-sizing: border-box;
    --svrollbar-thumb-width: 10px
  }
  /* Define the styles for the list of exchanges */
  .list {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .category {
    font-weight: bold;
    margin: 0.5rem 0;
    padding: 0.5rem;
    background-color: #202020;
  }
  /* Define the styles for the exchange items */
  .item {
    display: flex;
    align-items: center;
    margin: 0;
    padding: 0.4rem;
    width: 200px;
    cursor: pointer;
  }
  .tickbox {
    margin-right: 0.5rem;
  }
  .table {
    table-layout: fixed;
    width: 100%;
  }
  .th {
    color: white;
    font-weight: bold;
    padding: 0.5rem;
    border-bottom: 1px solid #d0d0d0;
  }
  .td1 {
    color: white;
    padding: 0.5rem;
    border-bottom: 1px solid #d0d0d0;
    overflow: hidden;
  }
  .td2 {
    color: white;
    padding: 0.5rem;
    border-bottom: 1px solid #d0d0d0;
    overflow: hidden;
  }
</style>
