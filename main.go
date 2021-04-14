package main

import (
	"github.com/pquerna/ffjson/ffjson"
	"path"
)
import (
	"GoCoinbaseFetcher/datastructure"
	//"encoding/json"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

const API = `https://api.pro.coinbase.com/products/%s/trades`

const BTC_FILE = `data/btc-eur_%s.json`
const ETH_FILE = `data/eth-eur_%s.json`
const LTC_FILE = `data/ltc-eur_%s.json`

func main() {

	var historyBTCTrades strings.Builder
	var historyETHTrades strings.Builder
	var historyLTCTrades strings.Builder

	historyBTCTrades.WriteString("[")
	historyETHTrades.WriteString("[")
	historyLTCTrades.WriteString("[")

	defer dumpAllData(historyBTCTrades, historyETHTrades, historyLTCTrades, "PANIC")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(dumpAllData(historyBTCTrades, historyETHTrades, historyLTCTrades, "CTRL+C"))
	}()

	var i int
	for {
		historyBTCTrades.WriteString(getHistoryString("BTC-EUR"))
		historyETHTrades.WriteString(getHistoryString("ETH-EUR"))
		historyLTCTrades.WriteString(getHistoryString("LTC-EUR"))
		log.Println(i)
		time.Sleep(2500 * time.Millisecond)
		i++
	}
}

func MergeData(target, finalName string) {
	files := fileutils.FindFiles("data", target, true)
	var data []datastructure.Trade
	for _, f := range files {
		if !strings.HasSuffix(f, ".json") {
			continue
		}
		log.Println("Managing file: " + f)
		var tempData []datastructure.Trade

		file, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		if err = ffjson.Unmarshal(file, &tempData); err != nil {
			panic(err)
		}
		data = append(data, tempData...)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].TradeID < data[j].TradeID
	})

	buf, err := ffjson.Marshal(data)
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
	ffjson.Pool(buf)
}

func dumpAllData(historyBTCTrades, historyETHTrades, historyLTCTrades strings.Builder, message string) int {
	log.Println(message + " INTERCEPTED!")
	timeNow := time.Now().Format("2006.01.02_15.04.05")

	// Dumping BTC
	dumpData(historyBTCTrades.String(), fmt.Sprintf(BTC_FILE, timeNow))
	debug.FreeOSMemory()
	historyBTCTrades.Reset()

	// DUMPING ETH
	dumpData(historyETHTrades.String(), fmt.Sprintf(ETH_FILE, timeNow))
	debug.FreeOSMemory()
	historyETHTrades.Reset()

	// DUMPING LTC
	dumpData(historyLTCTrades.String(), fmt.Sprintf(LTC_FILE, timeNow))
	debug.FreeOSMemory()
	historyLTCTrades.Reset()
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
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	rawData = rawData[1 : len(rawData)-1]
	return string(rawData) + ","
}
