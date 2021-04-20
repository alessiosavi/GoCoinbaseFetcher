package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"time"
)

const API_PRICE_HISTORY = `https://api.coinbase.com/v2/prices/BTC-EUR/spot?date=%s`

var client = http.Client{
	Transport:     http.DefaultTransport,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       0,
}

type TradeEntry struct {
	Amount   json.Number `json:"amount"`
	Base     string      `json:"base"`
	Currency string      `json:"currency"`
	Date     string      `json:"date"`
}

func FetchAllData(BTC_FILE_EUR string) {

	timeNow := time.Now().Format("2006.01.02_15.04.05")
	f, err := os.Create(fmt.Sprintf(BTC_FILE_EUR, timeNow))
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)
	w.WriteString("[")

	defer dumpAllData(f)
	defer f.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		dumpAllData(f)
		f.Close()
		os.Exit(0)
	}()

	date, err := time.Parse("2006-01-02", "2021-04-20")
	if err != nil {
		return
	}

	for {
		//if _, err = w.Write(fetchTradesHistory(&pagination, before)); err != nil {
		//	panic(err)
		//}
		if _, err = w.Write(fetchPricesHistory(&date)); err != nil {
			panic(err)
		}

		if err = w.Flush(); err != nil {
			panic(err)
		}

		log.Println("Date: ", date.Format("2006-01-02"))
		time.Sleep(400 * time.Millisecond)
	}
}

func dumpAllData(f *os.File) int {
	if _, err := f.Seek(-1, io.SeekEnd); err != nil {
		panic(err)
	}
	f.WriteString("]")
	return 0
}

func fetchPricesHistory(date *time.Time) []byte {
	resp, err := client.Get(fmt.Sprintf(API_PRICE_HISTORY, date.Format("2006-01-02")))
	if err != nil {
		log.Println("ERROR:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		response, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		log.Println("TOO MUCH REQUEST:\n" + string(response))
		time.Sleep(6 * time.Second)
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if !bytes.Contains(rawData, []byte("message")) && len(rawData) > 0 {
		var tradeEntry map[string]*TradeEntry = make(map[string]*TradeEntry)
		if err = json.Unmarshal(rawData, &tradeEntry); err != nil {
			log.Println("Unable to unmarshal data into tradeEntry:", err)
			return nil
		}
		tradeEntry["data"].Date = date.Format("2006-01-02")
		marshal, err := json.Marshal(tradeEntry["data"])
		if err != nil {
			log.Println("Unable to marshal data:", err)
			return nil
		}
		*date = date.Add(-24 * time.Hour)
		return append(marshal, []byte{','}...)
	}
	return nil
}
