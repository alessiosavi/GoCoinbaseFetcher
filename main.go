package main

import (
	"GoCoinbaseFetcher/datastructure"
	"encoding/json"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"github.com/alessiosavi/GoGPUtils/helper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"time"
)

const API = `https://api.pro.coinbase.com/products/%s/trades`

func main() {

	var historyBTCTrades []datastructure.Trade
	var historyETHTrades []datastructure.Trade
	var historyLTCTrades []datastructure.Trade

	defer dumpAllData(historyBTCTrades, historyETHTrades, historyLTCTrades, "PANIC")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(dumpAllData(historyBTCTrades, historyETHTrades, historyLTCTrades, "CTRL+C"))
	}()

	var i int
	for {
		historyBTCTrades = append(historyBTCTrades, getHistory("BTC-EUR")...)
		historyETHTrades = append(historyETHTrades, getHistory("ETH-EUR")...)
		historyLTCTrades = append(historyLTCTrades, getHistory("LTC-EUR")...)
		log.Println(i)
		time.Sleep(2500 * time.Millisecond)
		i++
	}
}

func dumpAllData(historyBTCTrades []datastructure.Trade, historyETHTrades []datastructure.Trade, historyLTCTrades []datastructure.Trade, message string) int {
	log.Println(message + " INTERCEPTED!")
	historyBTCTrades = append(loadData("data/btc-eur.json"), historyBTCTrades...)
	dumpData(historyBTCTrades, "data/btc-eur.json")
	historyBTCTrades = nil
	debug.FreeOSMemory()
	historyETHTrades = append(loadData("data/eth-eur.json"), historyETHTrades...)
	dumpData(historyETHTrades, "data/eth-eur.json")
	historyETHTrades = nil
	debug.FreeOSMemory()
	historyLTCTrades = append(loadData("data/ltc-eur.json"), historyLTCTrades...)
	dumpData(historyLTCTrades, "data/ltc-eur.json")
	historyLTCTrades = nil
	debug.FreeOSMemory()
	return 0
}

func dumpData(historyTrades []datastructure.Trade, filename string) {
	sort.Slice(historyTrades, func(i, j int) bool {
		return historyTrades[i].TradeID < historyTrades[j].TradeID
	})
	indent := helper.MarshalIndent(historyTrades)
	_ = ioutil.WriteFile(filename, []byte(indent), 0755)
}

func loadData(file string) []datastructure.Trade {
	var historyTrades []datastructure.Trade
	if fileutils.FileExists(file) {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if err = json.Unmarshal(data, &historyTrades); err != nil {
			panic(err)
		}
	}
	return historyTrades
}

func getHistory(pair string) []datastructure.Trade {
	resp, err := http.Get(fmt.Sprintf(API, pair))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 429 {
		response, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		panic("TOO MUCH REQUEST:\n" + string(response))
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var trades []datastructure.Trade
	if err = json.Unmarshal(rawData, &trades); err != nil {
		panic(err)
	}
	return trades
}
