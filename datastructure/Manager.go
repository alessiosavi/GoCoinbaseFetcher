package datastructure

import (
	"encoding/csv"
	"fmt"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Manager struct {
}

func (m Manager) DropDuplicates(filename string) {
	var history History
	history.Load(filename)
	history.DropDuplicates()
	csvdata := []byte(stringutils.JoinSeparator("|", HEADERS...) + "\n")
	csvdata = append(csvdata, []byte(history.CSV())...)
	ioutil.WriteFile(filename, csvdata, 0755)
}
func (m Manager) Sort(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.ReuseRecord = true
	csvReader.Comma = '|'
	csvReader.TrimLeadingSpace = true

	// Ignore header
	_, _ = csvReader.Read()
	lines, err := fileutils.CountLines(filename, "\n", 4096)
	if err != nil {
		panic(err)
	}

	var data = make(History, lines-1)
	log.Println(fmt.Sprintf("Sorting file: %s | Lines: %d", filename, lines))
	for line := 0; line < lines-1; line++ {
		row, err := csvReader.Read()
		if err == io.EOF {
			log.Println(fmt.Sprintf("Read %d rows in a total of %d", line, lines))
			break
		}
		if err != nil {
			panic(err)
		}
		data[line] = New(row)
	}
	data.Sort()

	csvdata := []byte(stringutils.JoinSeparator("|", HEADERS...) + "\n")
	csvdata = append(csvdata, []byte(data.CSV())...)
	ioutil.WriteFile(filename, csvdata, 0755)
}
