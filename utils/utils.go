package utils

import (
	"GoCoinbaseFetcher/datastructure"
	"bufio"
	"bytes"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http/httputil"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TradePagination struct {
	Name string `json:"name"`
	ID   string `json:"before"`
}

func FetchAllData(API, BTC_FILE_EUR, tradeID string, before bool) {
	var pagination TradePagination
	pagination.Name = fmt.Sprintf(API, "btc-eur")
	pagination.ID = tradeID

	log.Println("Using the following pagination: " + tradeID)

	timeNow := time.Now().Format("2006.01.02_15.04.05")
	f, err := os.Create(fmt.Sprintf(BTC_FILE_EUR, timeNow))
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)
	w.WriteString("[")

	defer DumpAllData(f)
	defer f.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		DumpAllData(f)
		f.Close()
		os.Exit(0)
	}()

	for {
		if _, err = w.Write(fetch(&pagination, before)); err != nil {
			panic(err)
		}
		if err = w.Flush(); err != nil {
			panic(err)
		}

		log.Println("TRADE ID: ", pagination.ID)
		if !before {
			time.Sleep(5000 * time.Millisecond)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func GetPagination(coin string, before bool) int {
	var tradeId int
	if before {
		tradeId = math.MaxInt32
	} else {
		tradeId = 0
	}

	t, err := fileutils.ListFile("data")
	if err != nil {
		return tradeId
	}
	var files []string
	for _, f := range t {
		if strings.HasSuffix(f, ".json") {
			files = append(files, f)
		}
	}
	sort.Strings(files)
	for _, f := range files {
		if strings.Contains(f, coin) && strings.HasSuffix(f, ".json") {
			open, err := os.Open(f)
			if err != nil {
				log.Println(err)
				continue
			}
			if before { // Read from the last character
				if _, err = open.Seek(-101, io.SeekEnd); err != nil {
					log.Println("Unable to seek: ", err)
					continue
				}
			} else { // Read from the first character
				if _, err = open.Seek(0, io.SeekStart); err != nil {
					log.Println("Unable to seek: ", err)
					continue
				}
			}
			b := make([]byte, 200)
			if _, err = open.Read(b); err != nil { // Read the first 120 byte
				log.Println("Unable to read the 120 byte: ", err)
				continue
			}
			row := string(b)
			if strings.Contains(row, "message") || strings.Contains(row, "Invalid") {
				continue
			}
			var startIndex, stopIndex int
			startIndex = strings.Index(row, `"trade_id":`) + len(`"trade_id":`)
			if before {
				// This hack have to be done cause i've changed the order of the struct in order to align byte
				// and reduce memory usage of struct on rapsberry Pi :'(
				stopIndex = strings.Index(row[startIndex:], ",") + startIndex
				if stopIndex < startIndex {
					stopIndex = strings.Index(row[startIndex:], "}") + startIndex
				}
			} else {
				stopIndex = strings.Index(row[startIndex:], "}") + startIndex
				stopIndex2 := strings.Index(row[startIndex:], ",") + startIndex
				if stopIndex2 < stopIndex && stopIndex2>startIndex{
					stopIndex = stopIndex2
				}
			}
			tradeIdTemp, err := strconv.Atoi(row[startIndex:stopIndex])
			if err != nil {
				log.Println("Unable to cast tradeID to int:", err)
				continue
			}
			if before {
				if tradeIdTemp < tradeId {
					tradeId = tradeIdTemp
				}
			} else {
				if tradeIdTemp > tradeId {
					tradeId = tradeIdTemp
				}
			}
		}
	}
	log.Println("Using the following tradeID:", tradeId)
	return tradeId
}

func DumpAllData(f *os.File) int {
	if _, err := f.Seek(-1, io.SeekEnd); err != nil {
		panic(err)
	}
	f.WriteString("]")
	return 0
}

func fetch(conf *TradePagination, before bool) []byte {
	var url string
	if before {
		url = conf.Name + fmt.Sprintf("?after=%s", conf.ID)
	} else {
		url = conf.Name + fmt.Sprintf("?before=%s", conf.ID)
	}
	resp, err := datastructure.Client.Get(url)
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

	if !bytes.Contains(rawData, []byte("message")) && len(rawData) > 0 {
		rawData = rawData[1 : len(rawData)-1]
		if before {
			conf.ID = resp.Header.Get("CB-BEFORE")
		} else {
			conf.ID = resp.Header.Get("CB-AFTER")
		}
		return append(rawData, []byte{','}...)
	}
	return nil
}
