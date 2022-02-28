package main

import (
	"github.com/alessiosavi/coinbase-fetcher/datastructure"
	"log"
	"time"
)

//curl --request GET --url 'https://api.exchange.coinbase.com/products/BTC-USD/candles?granularity=60&start=start&end=end' --header 'Accept: application/json'

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	var opts = datastructure.DownloadOpts{
		Granularity: datastructure.GRANULARITY_15_MINUTES,
		Pair:        "BTC-USD",
		LimitDate:   time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	download := opts.New()
	filename := download.Download(nil)

	var manager datastructure.Manager
	manager.Sort(filename)
	manager.DropDuplicates(filename)
}
