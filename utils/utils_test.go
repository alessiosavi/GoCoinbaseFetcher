package utils

import (
	"log"
	"testing"
)

const API = `https://api.pro.coinbase.com/products/%s/trades`
const BTC_FILE_EUR = `data/btc-eur_%s.json`

func TestGetPagination(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	GetPagination("btc-eur")
}