package datastructure

import (
	"encoding/csv"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HistoricRateRaw [][]float64

func (h HistoricRateRaw) CSV() string {
	var sb strings.Builder
	for i := range h {
		sb.WriteString(stringutils.JoinSeparator("|", []string{
			fmt.Sprintf("%d", int64(h[i][0])),
			fmt.Sprintf("%f", h[i][1]),
			fmt.Sprintf("%f", h[i][2]),
			fmt.Sprintf("%f", h[i][3]),
			fmt.Sprintf("%f", h[i][4]),
			fmt.Sprintf("%f", h[i][5])}...))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (h HistoricRateRaw) IsSorted() bool {
	return sort.SliceIsSorted(h, func(i, j int) bool {
		return h[i][0] < h[j][0]
	})
}

func (h HistoricRateRaw) IsMissing(granularity Granularity) bool {
	for i := 0; i < len(h)-1; i++ {
		if math.Abs(h[i][0]-h[i+1][0]) > float64(granularity) {
			return true
		}
	}
	return false
}

type History []HistoricRate

func (h *History) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.ReuseRecord = true
	csvReader.Comma = '|'
	csvReader.TrimLeadingSpace = true

	lines, err := fileutils.CountLines(filename, "\n", 4096)
	if err != nil {
		panic(err)
	}

	*h = make(History, lines-1)

	// Ignore header
	_, _ = csvReader.Read()
	var t HistoricRate

	log.Println(fmt.Sprintf("Loading file: %s | Lines: %d", filename, lines))
	for line := 0; line < lines-1; line++ {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		t.New(row)
		(*h)[line] = t
	}
	log.Println("File loaded!")
}

func (h History) Sort() {
	log.Println("Sorting ...")
	sort.Slice(h, func(i, j int) bool {
		return h[i].Time < h[j].Time
	})
	log.Println("Sorted!")
}
func (h History) IsSort() bool {
	return sort.SliceIsSorted(h, func(i, j int) bool {
		return h[i].Time < h[j].Time
	})
}

func (h History) CSV() string {
	var sb strings.Builder
	for i := 0; i < len(h); i++ {
		sb.WriteString(stringutils.JoinSeparator("|", h[i].CSV()...))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (h *History) DropDuplicates() {
	log.Println("Dropping duplicates ...")
	if !h.IsSort() {
		h.Sort()
	}
	for i := 0; i < len(*h)-1; i++ {
		if (*h)[i].Equal((*h)[i+1]) {
			*h = append((*h)[:i], (*h)[i+1:]...)
		}
	}
	log.Println("Duplicates dropped!")
}

type HistoricRate struct {
	Time   int64   `json:"time"`
	Low    float64 `json:"low,omitempty"`
	High   float64 `json:"high,omitempty"`
	Open   float64 `json:"open,omitempty"`
	Close  float64 `json:"close,omitempty"`
	Volume float64 `json:"volume,omitempty"`
}

func (h HistoricRate) CSV() []string {
	return []string{
		fmt.Sprintf("%d", h.Time),
		fmt.Sprintf("%f", h.Low),
		fmt.Sprintf("%f", h.High),
		fmt.Sprintf("%f", h.Open),
		fmt.Sprintf("%f", h.Close),
		fmt.Sprintf("%f", h.Volume)}
}
func (h *HistoricRate) New(data []string) {
	if unix, err := strconv.ParseFloat(data[0], 64); err == nil {
		h.Time = int64(unix)
	}
	if low, err := strconv.ParseFloat(data[1], 64); err == nil {
		h.Low = low
	}

	if high, err := strconv.ParseFloat(data[2], 64); err == nil {
		h.High = high
	}

	if open, err := strconv.ParseFloat(data[3], 64); err == nil {
		h.Open = open
	}

	if c, err := strconv.ParseFloat(data[4], 64); err == nil {
		h.Close = c
	}
	if volume, err := strconv.ParseFloat(data[5], 64); err == nil {
		h.Volume = volume
	}
}
func New(data []string) HistoricRate {
	var h HistoricRate
	h.New(data)
	return h
}
func (h HistoricRate) Equal(target HistoricRate) bool {
	return h.Time == target.Time &&
		h.Low == target.Low &&
		h.High == target.High &&
		h.Open == target.Open &&
		h.Close == target.Close &&
		h.Volume == target.Volume
}

var HEADERS = []string{"timestamp", "low", "high", "open", "close", "volume"}

var CLIENT *http.Client

func init() {
	CLIENT = http.DefaultClient
	CLIENT.Timeout = time.Second * 60
}
