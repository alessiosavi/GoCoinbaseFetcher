package main

import "testing"

func TestMergeData(t *testing.T) {
	MergeData("btc-eur", "btc-eur.json")
	MergeData("eth-eur", "eth-eur.json")
	MergeData("ltc-eur", "ltc-eur.json")
	MergeData("btc-usd", "btc-usd.json")
	MergeData("eth-usd", "eth-usd.json")
	MergeData("ltc-usd", "ltc-usd.json")
}
