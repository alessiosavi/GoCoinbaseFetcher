package datastructure

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

const BASE_URL = `https://api.exchange.coinbase.com/products/%s/candles?granularity=%d&start=%s&end=%s`
const LAYOUT = `2006-01-02T15:04:05Z`
const MAX_ROW_RETURNED = 300

type Granularity uint32

const (
	//60, 300, 900, 3600, 21600, 86400
	GRANULARITY_MINUTE     Granularity = 60
	GRANULARITY_5_MINUTES              = GRANULARITY_MINUTE * 5
	GRANULARITY_15_MINUTES             = GRANULARITY_MINUTE * 15
	GRANULARITY_HOUR                   = GRANULARITY_MINUTE * 60
	GRANULARITY_6_HOURS                = GRANULARITY_HOUR * 6
	GRANULARITY_DAY                    = GRANULARITY_HOUR * 24
)

type DownloadOpts struct {
	Granularity Granularity
	Pair        string
	LimitDate   time.Time
}

type Download struct {
	DownloadOpts
}

func (o DownloadOpts) New() Download {
	return Download{
		DownloadOpts: DownloadOpts{
			Granularity: o.Granularity,
			Pair:        o.Pair,
			LimitDate:   o.LimitDate,
		},
	}
}

func (d Download) GetURL(startDate, endDate time.Time) string {
	return fmt.Sprintf(BASE_URL, d.Pair, d.Granularity, Format(startDate), Format(endDate))
}

func (d Download) Sleep() {
	if helper.RandomInt(0, 10)%2 == 0 {
		time.Sleep(time.Millisecond * time.Duration(helper.RandomInt(0, 100)))
	}
}

func (d Download) Request(startDate, endDate time.Time) (HistoricRateRaw, error) {
	URL := d.GetURL(startDate, endDate)
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	d.Sleep()
	response, err := CLIENT.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 429 {
		log.Println("Timeout!!")
		time.Sleep(time.Second)
		return d.Request(startDate, endDate)
	}
	if response.StatusCode != 200 {
		var desc string
		if msg, err := ioutil.ReadAll(response.Body); err != nil {
			desc = string(msg)
		}
		return nil, fmt.Errorf("error using startDate: [%s] endDate: [%s] | Status Code: %d | Error: %s", startDate, endDate, response.StatusCode, desc)

	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var h HistoricRateRaw
	if err = json.Unmarshal(responseBody, &h); err != nil {
		panic(err)
	}

	if h.IsMissing(d.Granularity) {
		log.Println(fmt.Sprintf("Missing some date for the following period: [%s]-[%s]", Format(startDate), Format(endDate)))
	}
	return h, nil
}

func (d Download) Download(endDate *time.Time) string {
	if endDate == nil {
		x := time.Now()
		endDate = &x
	}
	startDate := GetNextDate(d.Granularity, *endDate)
	filename := fmt.Sprintf("%s-%d-%s-%s.csv", d.Pair, d.Granularity, Format(d.LimitDate), Format(*endDate))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	headers := []string{"timestamp", "low", "high", "open", "close", "volume"}
	w.WriteString(stringutils.JoinSeparator("|", headers...) + "\n")

	total := math.Ceil(startDate.Sub(d.LimitDate).Seconds() / float64(d.Granularity) / 300)
	bar := progressbar.Default(int64(total))
	defer bar.Close()
	log.Println("Start Downloading!")
	log.Println(fmt.Sprintf("Start Date: %s | End Date: %s | Pair: %s | Granularity: %d", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), d.Pair, d.Granularity))
	for endDate.After(d.LimitDate) {
		data, err := d.Request(startDate, *endDate)
		if err != nil {
			log.Println("PANIC!", err)
		}
		if len(data) != 0 {
			bar.Describe(UnixToTime(fmt.Sprintf("%d", int64(data[0][0]))))
			w.WriteString(data.CSV())
			w.Flush()
		}
		*endDate = startDate
		startDate = GetNextDate(d.Granularity, *endDate)
		bar.Add(1)
	}
	w.Flush()
	return filename
}

// GetNextDate is delegated to calculate the new date for which we have to retrieve the data
func GetNextDate(granularity Granularity, date time.Time) time.Time {
	minute := time.Second * time.Duration(granularity)
	return date.Add(-(minute * MAX_ROW_RETURNED))
}
func Format(date time.Time) string {
	return date.Format(LAYOUT)
}

func UnixToTime(unix string) string {
	n, err := strconv.Atoi(unix)
	if err != nil {
		panic(err)
	}
	return time.Unix(int64(n), 0).Format(LAYOUT)
}
