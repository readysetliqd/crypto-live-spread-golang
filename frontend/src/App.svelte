<script lang="ts">
  import { EventsOn } from '../wailsjs/runtime'
  import { onMount } from 'svelte'
  import { Svroller } from "svrollbar"
  import * as go from '../wailsjs/go/main/App.js'

  class Exchange {
    category: Category
    name: string
    checked: boolean = false
    pairs: string[] = []
    recommended: string[] = []
    askVolume: number = null
    ask: number = Number.MAX_VALUE
    bid: number = null
    bidVolume: number = null

    constructor(name: string, category: Category) {
      this.name = name
      this.category = category
    }
  }

  class Stat {
    value: number = null
    exchange: string = null    
  }

  type Category = "Spot" | "Futures" | "Hybrid DEX"

  const exchangeMap = new Map<string, Category>([
    ["Binance", "Spot"],
    ["Binance US", "Spot"],
    ["Bitget", "Spot"],
    ["Coinbase", "Spot"],
    ["Okx", "Spot"],
    ["Kraken", "Spot"],
    ["Bybit", "Spot"],
    ["Upbit", "Spot"],
    ["Kraken (Futures)", "Futures"],
    ["Binance (USD-M)", "Futures"],
    ["Binance (COIN-M)", "Futures"],
    ["Bybit (Futures)", "Futures"],
    ["Okx (Swaps)", "Futures"],
    ["Bitget (Futures)", "Futures"],
    ["DYDX", "Hybrid DEX"],
    ["HyperliquidX", "Hybrid DEX"],
  ])

  let exchanges: { [key: string]: Exchange } = {}
  exchangeMap.forEach((category: Category, name: string) => {
    let newExchange = new Exchange(name, category)
    exchanges[name] = newExchange
  })

  onMount(() => {
    EventsOn("spreadData", (exchange, bidVolume, bid, ask, askVolume) => {
      exchanges[exchange].bidVolume = bidVolume
      exchanges[exchange].bid = bid
      exchanges[exchange].ask = ask
      exchanges[exchange].askVolume = askVolume
    })
    EventsOn("HTTP Forbidden", () => {
      alert("HTTP Error 403 Forbidden\nYour IP may be geo-blocked by the server. Check exchange's supported locations and try again.")
    })
  })

  let coin: string = ''
  let denom: string = ''
  let highestBid = new Stat()
  let lowestAsk = new Stat()
  let widestSpread = new Stat()
  let highLowDiff = new Stat()
  let exchangeNames = Object.keys(exchanges)

  $: { exchanges,
    highestBid.value = 0
    lowestAsk.value = Number.MAX_VALUE
    widestSpread.value = 0
    Object.keys(exchanges).forEach((exchange) => {
      let spread = 0
      if (exchanges[exchange].bid != null) {
        spread = (exchanges[exchange].ask / exchanges[exchange].bid * 100 - 100)
      }
      if (spread > widestSpread.value) {
        widestSpread.value = spread
        widestSpread.exchange = exchange
      }
      if (exchanges[exchange].bid > highestBid.value) {
        highestBid.value = exchanges[exchange].bid
        highestBid.exchange = exchange
      } 
      if (exchanges[exchange].ask < lowestAsk.value) {
        lowestAsk.value = exchanges[exchange].ask
        lowestAsk.exchange = exchange
      }
    })
  }

  $ : {lowestAsk.value, highestBid.value,
    highLowDiff.value = (lowestAsk.value / highestBid.value * 100 - 100)
  }

  $: {coin, denom, 
    Object.keys(exchanges).forEach((exchange) => {
      exchanges[exchange].recommended = getRecommendedPairs(exchange)
    })
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
    if (exchange == "HyperliquidX") {
      return exchanges[exchange].pairs.filter(pair => pair.includes(coinInput))
    } else {
      return exchanges[exchange].pairs.filter(pair => pair.includes(coinInput) && pair.includes(denomInput))
    }
  }

  async function fetchPairs(exchange): Promise<void> {
    console.log('fetchPairs called', exchange)
    const result = await go.FetchPairs(exchange)
    exchanges[exchange].pairs = result
  }
  
  async function connectWebsocket(exchange, pair): Promise<void> {
    console.log('connectWebsocket called', exchange)
    await go.ConnectWebsocket(exchange, pair)
  }

  const selectPair = (event) => {
    const [selectedPair, exchange] = event.target.value.split(',')
    connectWebsocket(exchange, selectedPair)
  }

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
  <table class="table">
    <thead>
      <tr class="stat-tr">
        <th></th>
        <th>Highest Bid</th>
        <th>Lowest Ask</th>
        <th>Widest Spread</th>
        <th>Highest/Lowest diff</th>
      </tr>
      <tr class="stat-tr">
        <td>Value</td>
        <td class="stat-td">{highestBid.value}</td>
        <td class="stat-td">{lowestAsk.value}</td>
        <td class="stat-td">{widestSpread.value.toFixed(3)}%</td>
        <td class="stat-td">{highLowDiff.value.toFixed(3)}%</td>
      </tr>
      <tr class="stat-tr">
        <td>Exchange</td>
        <td class="stat-td">{highestBid.exchange}</td>
        <td class="stat-td">{lowestAsk.exchange}</td>
        <td class="stat-td">{widestSpread.exchange}</td>
        <td class="stat-td">{highLowDiff.exchange}</td>
      </tr>
    </thead>
  </table>
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
  .stat-tr {
    color: white;
    border-bottom: 1px solid white;
  }
  .stat-td {
    color: white;
    border: 1px solid white;
  }
</style>
