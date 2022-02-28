package datastructure

import (
	"encoding/json"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data := `[[{"time":1637712000,"low":57562.65,"high":57607.97,"open":57562.65,"close":57599.99,"volume":7.355625},{"time":1637715480,"low":57601.58,"high":57601.58,"open":57601.58,"close":57601.58,"volume":0.000087}],[{"time":1632706200,"low":43914.5,"high":43953.16,"open":43939.67,"close":43918.1,"volume":6.544966},{"time":1632706320,"low":43922.66,"high":43944.87,"open":43922.67,"close":43937.35,"volume":6.614559}]]`
	var h [][]HistoricRate

	if err := json.Unmarshal([]byte(data), &h); err != nil {
		panic(err)
	}
	if len(h) != 2 {
		log.Println(helper.MarshalIndent(h))
		panic(len(h))
	}

}
