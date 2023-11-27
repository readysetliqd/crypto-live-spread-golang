<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime'
  import {FetchKrakenSpotPairs} from '../wailsjs/go/main/App.js'
  import {FetchBinanceSpotPairs} from '../wailsjs/go/main/App.js'
  import {ConnectKrakenSpotWebsocket} from '../wailsjs/go/main/App.js'
  import { Svroller } from "svrollbar"

  let coin: string
  let denom: string
  let selectedAsset: string = "";
  let krakenSpread: number[] = [1, 2, 3]
  let pairsPlaceholder: string[] = []

 function fetchPairs(exchange) {
  console.log('fetchPairs called', exchange)
  switch (exchange) {
    case 'Kraken':
    fetchKrakenSpotPairs()
    case 'Binance':
    fetchBinanceSpotPairs()
  }
 }

  async function fetchKrakenSpotPairs(): Promise<void> {
    console.log('fetchKrakenSpotPairs called');
    const result = await FetchKrakenSpotPairs();
    console.log('FetchKrakenSpotPairs result', result);
    exchanges["Kraken"].pairs = result;
  }

  async function fetchBinanceSpotPairs(): Promise<void> {
    console.log('fetchBinanceSpotPairs called')
    const result = await FetchBinanceSpotPairs();
    console.log('FetchBinanceSpotPairs result', result);
    exchanges["Binance"].pairs = result
  }

  async function connectKrakenSpotWebsocket(asset): Promise<void> {
    console.log('connectKrakenSpotWebsocket called');
    const result = await ConnectKrakenSpotWebsocket(asset)
    console.log('ConnectKrakenSpotWebsocket result', result)
  }

  onMount(() => {
    EventsOn("spreadData", (bidVolume, bid, ask, askVolume) => {
      krakenSpread = [bidVolume, bid, ask, askVolume];
    });
  });

  const selectAsset = (event) => {
    selectedAsset = event.target.value;
    connectKrakenSpotWebsocket(selectedAsset)
  }

  let exchanges = {
    "Kraken": {category: "Spot", name: "Kraken", checked: false, pairs: pairsPlaceholder},
    "Binance": {category: "Spot", name: "Binance", checked: false, pairs: pairsPlaceholder},
    "Binance US": {category: "Spot", name: "Binance US", checked: false, pairs: pairsPlaceholder},
    "Okx": {category: "Spot", name: "Okx", checked: false, pairs: pairsPlaceholder},
    "Bitget": {category: "Spot", name: "Bitget", checked: false, pairs: pairsPlaceholder},
    "Bybit": {category: "Spot", name: "Bybit", checked: false, pairs: pairsPlaceholder},
    "Upbit": {category: "Spot", name: "Upbit", checked: false, pairs: pairsPlaceholder},
    "Kraken (Futures)": {category: "Futures", name: "Kraken (Futures)", checked: false, pairs: pairsPlaceholder},
    "Binance (USD-M)": {category: "Futures", name: "Binance (USD-M)", checked: false, pairs: pairsPlaceholder},
    "Binance (COIN-M)": {category: "Futures", name: "Binance (COIN-M)", checked: false, pairs: pairsPlaceholder},
    "Bybit (Futures)": {category: "Futures", name: "Bybit (Futures)", checked: false, pairs: pairsPlaceholder},
    "Okx (Futures)": {category: "Futures", name: "Okx (Futures)", checked: false, pairs: pairsPlaceholder},
    "Bitget (Futures)": {category: "Futures", name: "Bitget (Futures)", checked: false, pairs: pairsPlaceholder},
    "HyperliquidX": {category: "DEX", name: "HyperliquidX", checked: false, pairs: pairsPlaceholder},
    "DYDX": {category: "DEX", name: "DYDX", checked: false, pairs: pairsPlaceholder},
    "GMX": {category: "DEX", name: "GMX", checked: false, pairs: pairsPlaceholder},
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
        <li class="item" on:click={() => toggleExchange(exchanges[exchangeName])}>
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
  <div class="body">
    {#each exchangeNames as exchange}
    {#if exchanges[exchange].checked}
    <table class="table">
      <thead>
        <tr>

          <th class="th">{exchange}</th>
          <th class="th">Bid Volume</th>
          <th class="th">Ask Price</th>
          <th class="th">Ask Price</th>
          <th class="th">Ask Volume</th>

        </tr>
      </thead>
      <tbody>
        <tr>

          <td class="td">
            <select class="dropdown">
              {#each exchanges[exchange].pairs as pair}
              <option value="{pair}">{pair}</option>
              {/each}
              </select>
          </td>
          <td class="td">0</td>
          <td class="td">0</td>
          <td class="td">0</td>
          <td class="td">0</td>

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
  .title {
    height: 20px;
    line-height: 20px;
    margin-top: 1.5rem;
    margin-bottom: 0.5rem;
    margin-left: 12rem;
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
    margin-left: 12rem;
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
  /* Define the styles for the category headings */
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
  /* Define the styles for the tickboxes */
  .tickbox {
    margin-right: 0.5rem;
  }
    /* Define the styles for the table headings */
  .th {
    color: white;
    font-weight: bold;
    padding: 0.5rem;
    border-bottom: 1px solid #d0d0d0;
  }
  /* Define the styles for the table cells */
  .td {
    color: white;
    padding: 0.5rem;
    border-bottom: 1px solid #d0d0d0;
  }
  .table {
    margin-left:13rem
  }

</style>
