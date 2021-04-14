package main

import "testing"

func TestMergeData(t *testing.T) {
	MergeData("btc", "btc-eur.json")
	MergeData("eth", "eth-eur.json")
	MergeData("ltc", "ltc-eur.json")
}
