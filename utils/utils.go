package utils

import (
	"io"
	"log"
	"os"
)

type TradePagination struct {
	Name string `json:"name"`
	ID   string `json:"before"`
}

func DumpAllData(f *os.File) int {
	log.Println("Dumping data ...")
	if _, err := f.Seek(-1, io.SeekEnd); err != nil {
		panic(err)
	}
	f.WriteString("]")
	return 0
}
