package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"
)

type TradePagination struct {
	Name   string `json:"name"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func FetchAllData(API, BTC_FILE_EUR string) {
	var sb strings.Builder
	var pagination TradePagination
	pagination.Name = fmt.Sprintf(API, "btc-eur")
	pagination.After = fmt.Sprintf("%d", math.MaxInt32)

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
