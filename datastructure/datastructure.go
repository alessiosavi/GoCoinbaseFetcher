package datastructure

import "time"

type Trade struct {
	Time    time.Time `json:"time"`
	TradeID int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Side    string    `json:"side"`
}
