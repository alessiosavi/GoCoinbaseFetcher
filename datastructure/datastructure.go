package datastructure

import (
	"encoding/json"
	"time"
)

type Trade struct {
	Time    time.Time `json:"time"`
	TradeID int       `json:"trade_id"`
	Price   json.Number   `json:"price"`
	Size    json.Number   `json:"size"`
	Side    string    `json:"side"`
}
