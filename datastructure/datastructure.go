package datastructure

import (
	"encoding/json"
	"net/http"
	"time"
)

type Trade struct {
	Time    time.Time   `json:"time"`
	Price   json.Number `json:"price"`
	Size    json.Number `json:"size"`
	Side    string      `json:"side"`
	TradeID int         `json:"trade_id"`
}

type Candle struct{

}

var Client = http.Client{
	Transport:     http.DefaultTransport,
	Timeout:       30*time.Second,
}

