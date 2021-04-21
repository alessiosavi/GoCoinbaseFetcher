package api

import (
	"GoCoinbaseFetcher/datastructure"
	"GoCoinbaseFetcher/utils"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"time"
)

const CANDLE_API = "https://api-public.sandbox.pro.coinbase.com/products/%s/candles?start=%s&graularity=%d"
const CANDLE_FILE = "data/candle/%s/candle_%s_%s.json"

func GetCandles(pair, startDate string, granularity int) error {
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(path.Base(CANDLE_FILE), 0755); err != nil {
		return err
	}

	timeNow := time.Now().Format("2006.01.02_15.04.05")

	fName := fmt.Sprintf(CANDLE_FILE, pair, pair, timeNow)

	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer utils.DumpAllData(f)
	defer f.Close()
	w := bufio.NewWriter(f)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	t := time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)
	for date.After(t) {

		log.Printf("Date: %s\n", date)
		url := fmt.Sprintf(CANDLE_API, pair, date.Format("2006-01-02"), granularity)
		get, err := datastructure.Client.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}
		defer get.Body.Close()
		data, err := ioutil.ReadAll(get.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		if _, err = w.Write(append(data, []byte(",")...)); err != nil {
			log.Println(err)
			continue
		}
		if err = w.Flush(); err != nil {
			log.Println(err)
			continue
		}

		date = date.Add(-24 * time.Hour)
		time.Sleep(1 * time.Second)
	}
	return nil

}
