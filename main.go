package main

import (
	"GoCoinbaseFetcher/utils"
	"encoding/json"
	"flag"
	"path"
)
import (
	"GoCoinbaseFetcher/datastructure"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"time"
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
	b := flag.Bool("merge", false, "merge data")
	flag.Parse()
	if *b {
		MergeData("btc-eur", "btc-eur.json")
		return
	}
	utils.FetchAllData(API, BTC_FILE_EUR, fmt.Sprintf("%d", utils.GetPagination("btc-eur")))
}

func MergeData(target, finalName string) {
	files := fileutils.FindFiles("data", target, true)
	var data []datastructure.Trade
	for _, file := range files {
		if !strings.HasSuffix(file, ".json") {
			continue
		}
		log.Println("Managing file: " + file)

		open, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer open.Close()
		decoder := json.NewDecoder(open)

		if _, err = decoder.Token(); err != nil {
			panic(err)
		}
		for decoder.More() {
			var tempData datastructure.Trade
			if err = decoder.Decode(&tempData); err != nil {
				panic(err)
			}
			data = append(data, tempData)
		}
		open.Close()
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].TradeID < data[j].TradeID
	})

	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fName := path.Join("data/", finalName)
	if err = ioutil.WriteFile(fName, buf, 0755); err != nil {
		panic(err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f, ".json") || strings.Contains(f, finalName) {
			continue
		}
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func dumpAllData(historyBTCTrades, historyBTCTradesUSD, historyETHTradesUSD, historyETHTrades, historyLTCTrades, historyLTCTradesUSD strings.Builder, message string) int {
	log.Println(message + " INTERCEPTED!")
	timeNow := time.Now().Format("2006.01.02_15.04.05")

	// Dumping BTC
	dumpData(historyBTCTrades.String(), fmt.Sprintf(BTC_FILE_EUR, timeNow))
	historyBTCTrades.Reset()
	debug.FreeOSMemory()

	dumpData(historyBTCTradesUSD.String(), fmt.Sprintf(BTC_FILE_USD, timeNow))
	historyBTCTradesUSD.Reset()
	debug.FreeOSMemory()

	// DUMPING ETH
	dumpData(historyETHTrades.String(), fmt.Sprintf(ETH_FILE_EUR, timeNow))
	historyETHTrades.Reset()
	debug.FreeOSMemory()

	dumpData(historyETHTradesUSD.String(), fmt.Sprintf(ETH_FILE_USD, timeNow))
	historyETHTradesUSD.Reset()
	debug.FreeOSMemory()

	// DUMPING LTC
	dumpData(historyLTCTrades.String(), fmt.Sprintf(LTC_FILE_EUR, timeNow))
	historyLTCTrades.Reset()
	debug.FreeOSMemory()

	dumpData(historyLTCTradesUSD.String(), fmt.Sprintf(LTC_FILE_USD, timeNow))
	historyLTCTradesUSD.Reset()
	debug.FreeOSMemory()

	return 0
}

func dumpData(historyTrades, filename string) {
	_ = ioutil.WriteFile(filename, []byte(historyTrades[:len(historyTrades)-1]+"]"), 0755)
}

func getHistoryString(pair string) string {
	resp, err := http.Get(fmt.Sprintf(API, pair))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 429 {
		response, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		log.Println("TOO MUCH REQUEST:\n" + string(response))
		time.Sleep(5 * time.Second)
	}
	resp.Header.Get("after")
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	rawData = rawData[1 : len(rawData)-1]
	return string(rawData) + ","
}
