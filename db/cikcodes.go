package db

import (
	"bufio"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const cikURL = "https://www.sec.gov/include/ticker.txt"

var ccm *CikCodeMap

func init() {
	ccm = NewCikCodeMap()
	ccm.Build()
}

// CikCodeMap maps an equity symbol to a SEC CIK code
type CikCodeMap struct {
	urlEndpoint string
	codes       map[string]int
}

// NewCikCodeMap initializes a CikCodeMap
func NewCikCodeMap() *CikCodeMap {
	var ccm CikCodeMap
	ccm.urlEndpoint = cikURL
	ccm.codes = make(map[string]int)
	return &ccm
}

// Build reaches out to the sec url endpoint to parse
// and populate the CikCodeMap
func (ccm *CikCodeMap) Build() error {
	resp, err := http.Get(ccm.urlEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		symb, cik := parsecikcode(line)
		ccm.codes[symb] = cik
	}
	return nil
}

// Find searches the map to return the CIK code for the given symbol
func (ccm *CikCodeMap) Find(symbol string) int {
	return ccm.codes[symbol]
}

func parsecikcode(bs []byte) (symbol string, cik int) {
	d := strings.Fields(string(bs))
	s := strings.ToUpper(d[0])
	i, err := strconv.Atoi(d[1])
	if err != nil {
		panic(err)
	}
	return s, i
}
