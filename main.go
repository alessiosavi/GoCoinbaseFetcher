package main

import (
	"GoCoinbaseFetcher/datastructure"
	"GoCoinbaseFetcher/utils"
	"encoding/json"
	"flag"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"log"
	"os"
	"path"
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

	if !fileutils.IsDir("log") {
		os.Mkdir("log", 0755)
	}
	//
	//f, err := os.OpenFile(fmt.Sprintf("log/log_%s.log", time.Now().Format("2006.01.02_15.04.05")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//log.SetOutput(f)

	b := flag.Bool("merge", false, "merge data")
	mode := flag.Bool("before", false, "use `true` or `false` in order to download the past or future transaction")
	flag.Parse()
	if *b {
		log.Println("Merging data")
		MergeData("btc-eur", "btc-eur.json")
		return
	}

	if *mode {
		log.Println("Downloading old data ...")
	} else {
		log.Println("Downloading new data ...")
	}
	utils.FetchAllData(API, BTC_FILE_EUR, fmt.Sprintf("%d", utils.GetPagination("btc-eur", *mode)), *mode)
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
	log.Println("File read, going to sort ...")
	sort.Slice(data, func(i, j int) bool {
		return data[i].TradeID < data[j].TradeID
	})

	fName := path.Join("data/", finalName)
	log.Println("File Sorted! Going to dump into: " + fName)

	// If the file doesn't exist, create it
	f, err := os.OpenFile(fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(data); err != nil {
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file, ".json") || strings.Contains(file, finalName) {
			continue
		}
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}
}
