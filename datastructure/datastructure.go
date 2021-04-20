package datastructure

import (
	"encoding/json"
	"time"
)

type Trade struct {
	Time    time.Time   `json:"time"`
	Price   json.Number `json:"price"`
	Size    json.Number `json:"size"`
	Side    string      `json:"side"`
	TradeID int         `json:"trade_id"`
}
