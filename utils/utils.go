package utils

import (
	"bufio"
	"bytes"
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

var client = http.Client{
	Transport:     http.DefaultTransport,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       0,
}

func FetchAllData(API, BTC_FILE_EUR, tradeID string) {
	var pagination TradePagination
	pagination.Name = fmt.Sprintf(API, "btc-eur")
	pagination.After = tradeID

	log.Println("Using the following pagination: " + tradeID)

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

	for {
		if _, err = w.Write(fetch(&pagination)); err != nil {
			panic(err)
		}
		if err = w.Flush(); err != nil {
			panic(err)
		}

		log.Println("After: ", pagination.After)
		time.Sleep(500 * time.Millisecond)
	}
}

func GetPagination(coin string) int {
	t, err := fileutils.ListFile("data")
	if err != nil {
		return math.MaxInt32
	}

	var files []string
	for _, f := range t {
		if strings.HasSuffix(f, ".json") {
			files = append(files, f)
		}
	}
	sort.Strings(files)
	var tradeId int = math.MaxInt32
	for _, f := range files {
		if strings.Contains(f, coin) && strings.HasSuffix(f, ".json") {
			open, err := os.Open(f)
			if err != nil {
				panic(err)
			}
			if _, err = open.Seek(-101, io.SeekEnd); err != nil {
				log.Println("Unable to seek!")
				continue
			}
			b := make([]byte, 120)
			if _, err = open.Read(b); err != nil {
				panic(err)
			}
			row := string(b)
			startIndex := strings.Index(row, `"trade_id":`) + len(`"trade_id":`)
			stopIndex := strings.Index(row[startIndex:], ",") + startIndex
			tradeIdTemp, err := strconv.Atoi(row[startIndex:stopIndex])
			if err != nil {
				continue
			}
			if tradeIdTemp < tradeId {
				tradeId = tradeIdTemp
			}
		}
	}
	return tradeId
}

func dumpAllData(f *os.File) int {
	if _, err := f.Seek(-1, io.SeekEnd); err != nil {
		panic(err)
	}
	f.WriteString("]")
	return 0
}

func fetch(conf *TradePagination) []byte {
	url := conf.Name + fmt.Sprintf("?after=%s", conf.After)
	resp, err := client.Get(url)
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
		time.Sleep(5 * time.Second)
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if !bytes.Contains(rawData, []byte("error")) {
		rawData = rawData[1 : len(rawData)-1]
		conf.After = resp.Header.Get("CB-AFTER")
		return append(rawData, []byte{','}...)
	}
	return nil
}
