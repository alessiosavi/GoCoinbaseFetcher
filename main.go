package main

import (
	"GoCoinbaseFetcher/api"
	"log"
)

const API = `https://api.pro.coinbase.com/products/%s/trades`

const BTC_FILE_EUR = `data/btc-eur_%s.json`
const ETH_FILE_EUR = `data/eth-eur_%s.json`
const LTC_FILE_EUR = `data/ltc-eur_%s.json`

const BTC_FILE_USD = `data/btc-usd_%s.json`
const ETH_FILE_USD = `data/eth-usd_%s.json`
const LTC_FILE_USD = `data/ltc-usd_%s.json`

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	api.GetCandles("BTC-EUR", "2021-04-22", 60)
}
