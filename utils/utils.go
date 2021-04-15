package utils

import (
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TradePagination struct {
	Name   string `json:"name"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func FetchAllData(API, BTC_FILE_EUR, tradeID string) {
	var sb strings.Builder
	var pagination TradePagination
	pagination.Name = fmt.Sprintf(API, "btc-eur")
	pagination.After = tradeID

	log.Println("Using the following pagination: " + tradeID)
	sb.WriteString("[")

	defer dumpAllData(sb, BTC_FILE_EUR, "PANIC")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(dumpAllData(sb, BTC_FILE_EUR, "CTRL+C"))
	}()

	for {
		sb.WriteString(fetch(&pagination))
		log.Println("After: ", pagination.After)
		time.Sleep(500 * time.Millisecond)
	}
}

func GetPagination(coin string) int {
	t, err := fileutils.ListFile("C:\\opt\\SP\\workspace\\Golang\\GoCoinbaseFetcher\\data")
	if err != nil {
		panic(err)
	}

	var files []string
	for _, f := range t {
		if strings.HasSuffix(f, ".json") {
			files = append(files, f)
		}
	}

	sort.Strings(files)
	var tradeId int
	for _, f := range files {
		if strings.Contains(f, coin) {
			open, err := os.Open(f)
			if err != nil {
				panic(err)
			}
			if _, err = open.Seek(-120, io.SeekEnd); err != nil {
				panic(err)
			}

			b := make([]byte, 120)
			if _, err = open.Read(b); err != nil {
				panic(err)
			}
			row := string(b)
			log.Println(row)
			startIndex := strings.Index(row, `"trade_id":`) + len(`"trade_id":`)
			stopIndex := strings.Index(row[startIndex:], ",") + startIndex
			tradeId, err = strconv.Atoi(row[startIndex:stopIndex])
			if err != nil {
				return math.MaxInt32
			}
		}
	}
	return tradeId
}

func dumpAllData(data strings.Builder, BTC_FILE_EUR, message string) int {
	log.Println(message + " INTERCEPTED!")
	timeNow := time.Now().Format("2006.01.02_15.04.05")

	// Dumping BTC
	dumpData(data.String(), fmt.Sprintf(BTC_FILE_EUR, timeNow))
	data.Reset()
	debug.FreeOSMemory()
	return 0
}

func dumpData(historyTrades, filename string) {
	_ = ioutil.WriteFile(filename, []byte(historyTrades[:len(historyTrades)-1]+"]"), 0755)
}
func fetch(conf *TradePagination) string {
	url := conf.Name + fmt.Sprintf("?after=%s", conf.After)
	resp, err := http.Get(url)
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
	conf.After = resp.Header.Get("CB-AFTER")
	return string(rawData) + ","
}
