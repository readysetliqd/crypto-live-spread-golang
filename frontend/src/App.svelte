<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime'
  import {Greet} from '../wailsjs/go/main/App.js'
  import {FetchKrakenSpotPairs} from '../wailsjs/go/main/App.js'
  import {ConnectKrakenSpotWebsocket} from '../wailsjs/go/main/App.js'
  import { Svroller } from "svrollbar"

  let coin: string
  let denom: string
  let assets: {[key: string]: string;} = {}
  let selectedAsset: string = "";
  let krakenSpread: number[] = [1, 2, 3]

  async function fetchKrakenSpotPairs(): Promise<void> {
    console.log('fetchKrakenSpotPairs called');
    const result = await FetchKrakenSpotPairs();
    console.log('FetchKrakenSpotPairs result', result);
    assets = result;
  }

  async function connectKrakenSpotWebsocket(asset): Promise<void> {
    console.log('connectKrakenSpotWebsocket called');
    const result = await ConnectKrakenSpotWebsocket(asset)
    console.log('ConnectKrakenSpotWebsocket result', result)
  }

  onMount(() => {
    fetchKrakenSpotPairs();
    EventsOn("spreadData", (bidVolume, bid, ask, askVolume) => {
      krakenSpread = [bidVolume, bid, ask, askVolume];
    });
  });

  const selectAsset = (event) => {
    selectedAsset = event.target.value;
    connectKrakenSpotWebsocket(selectedAsset)
  }

  let exchanges = [
    { category: "Spot", name: "Kraken", checked: false },
    { category: "Spot", name: "Coinbase", checked: false },
    { category: "Spot", name: "Binance", checked: false },
    { category: "Spot", name: "Binance US", checked: false },
    { category: "Spot", name: "Okx", checked: false },
    { category: "Spot", name: "Bitget", checked: false },
    { category: "Spot", name: "Bybit", checked: false },
    { category: "Spot", name: "Upbit", checked: false },
    { category: "Futures", name: "Kraken (Futures)", checked: false },
    { category: "Futures", name: "Binance (USD-M)", checked: false },
    { category: "Futures", name: "Binance (COIN-M)", checked: false },
    { category: "Futures", name: "Bybit (Futures)", checked: false },
    { category: "Futures", name: "Okx (Futures)", checked: false },
    { category: "Futures", name: "Bitget (Futures)", checked: false },
    { category: "DEX", name: "HyperliquidX", checked: false },
    { category: "DEX", name: "DYDX", checked: false },
    { category: "DEX", name: "GMX", checked: false }
  ];
  // Define a function to toggle the checked status of an exchange
  function toggleExchange(exchange) {
    console.log("exchange checked", exchanges)
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
        {#each exchanges as exchange, i (exchange.name)}
        {#if i === 0 || exchange.category !== exchanges[i - 1].category}
        <li class="category">{exchange.category}</li>
        {/if}
        <li class="item" on:click={() => toggleExchange(exchange)}>
          <input
            class="tickbox"
            type="checkbox"
            bind:checked={exchange.checked}
          />
          {exchange.name}
        </li>
        {/each}
      </ul>
    </Svroller>
  </div>
  <div class="body">
    {#each exchanges as exchange (exchange.name)}
    {#if exchange.checked}
    <table class="table">
      <thead>
        <tr>

          <th class="th">{exchange.name}</th>
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
              {#each Object.entries(assets) as [_, value]}
              <option value="{value}">{value}</option>
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
  

  <div>
    <select bind:value={selectedAsset} on:change={selectAsset}>
      {#each Object.entries(assets) as [_, value]}
        <option value="{value}">{value}</option>
      {/each}
    </select>
  </div>
  <div class="selected-asset" id="selected-asset">{selectedAsset}</div>
  <div class="kraken-spread" id="kraken-spread">{krakenSpread}</div>
</main>

<style>

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
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
