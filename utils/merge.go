package utils

import (
	"GoCoinbaseFetcher/datastructure"
	"encoding/json"
	"github.com/alessiosavi/GoGPUtils/files"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

func MergeData(target, finalName string) {
	files := fileutils.FindFiles("data", target, true)
	var data []datastructure.Trade
	for _, file := range files {
		if !strings.HasSuffix(file, ".json") {
			continue
		}
		log.Println("Managing file: " + file)

		open, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer open.Close()
		decoder := json.NewDecoder(open)

		if _, err = decoder.Token(); err != nil {
			panic(err)
		}
		for decoder.More() {
			var tempData datastructure.Trade
			if err = decoder.Decode(&tempData); err != nil {
				panic(err)
			}
			data = append(data, tempData)
		}
		open.Close()
	}
	log.Println("File read, going to sort ...")
	sort.Slice(data, func(i, j int) bool {
		return data[i].TradeID < data[j].TradeID
	})

	fName := path.Join("data/", finalName)
	log.Println("File Sorted! Going to dump into: " + fName)

	// If the file doesn't exist, create it
	f, err := os.OpenFile(fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(data); err != nil {
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file, ".json") || strings.Contains(file, finalName) {
			continue
		}
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}
}
