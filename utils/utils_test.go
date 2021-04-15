package utils

import "testing"

const API = `https://api.pro.coinbase.com/products/%s/trades`
const BTC_FILE_EUR = `data/btc-eur_%s.json`

func TestFetchAllData(t *testing.T) {
	FetchAllData(API, BTC_FILE_EUR)
}
