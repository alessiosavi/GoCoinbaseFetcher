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
	"os"
	"sort"
	"strings"
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