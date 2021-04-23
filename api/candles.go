package api

import (
	"GoCoinbaseFetcher/utils"
	"bufio"
	"encoding/json"
	"fmt"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"log"
	"os"
	"os/signal"
	"path"
	"time"
)

const CANDLE_FILE = "data/candle/%s/candle_%s_%s.json"

var LAYOUT = "2006-01-02T15:04:05Z"

func GetHistoryCandles(pair, startDate string, granularity int) []byte {
	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL: "https://api.pro.coinbase.com",
	})
	client.RetryCount = 5
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Println(err)
		return nil
	}
	if err = os.MkdirAll(path.Base(CANDLE_FILE), 0755); err != nil {
		log.Println(err)
		return nil
	}

	timeNow := time.Now().Format("2006.01.02_15.04.05")

	fName := fmt.Sprintf(CANDLE_FILE, pair, pair, timeNow)

	f, err := os.Create(fName)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer utils.DumpAllData(f)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString("[")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		utils.DumpAllData(f)
		f.Close()
		os.Exit(0)
	}()

	t := time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)
	for date.After(t) {
		dateYesterday := date.Add(-5 * time.Hour)
		log.Printf("From: %s To: %s \n", date.UTC().Format(LAYOUT), dateYesterday.UTC().Format(LAYOUT))
		rates, err := client.GetHistoricRates(pair, coinbasepro.GetHistoricRatesParams{
			Start:       dateYesterday,
			End:         date,
			Granularity: granularity,
		})
		if err != nil {
			log.Println(err)
			log.Printf("ERROR! From: %+v To: %+v | ERR:\n%s\n", date.UTC().Format(LAYOUT), dateYesterday.UTC().Format(LAYOUT), err.Error())
			continue
		}
		date = dateYesterday

		data, err := json.Marshal(rates)
		if err != nil {
			log.Println(err)
			return nil
		}

		if _, err = w.Write(append(data, []byte(",")...)); err != nil {
			log.Println(err)
			continue
		}
		if err = w.Flush(); err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(400 * time.Millisecond)
	}
	return nil

}
