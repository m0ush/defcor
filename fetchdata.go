package fetcher

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	urlString   = "https://www.sec.gov/include/ticker.txt"
	csvfilename = "allsecurities.csv"
)

// CikCodeMap maps an equity symbol to a SEC CIK code
type CikCodeMap struct {
	urlEndpoint string
	codes       map[string]int
}

// SecuritySlice creates a slice of security symbols
type SecuritySlice struct {
	filename string
	symbols  []string
}

// NewSecuritySlice instatiates the security slice
func NewSecuritySlice(filename string) *SecuritySlice {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}
	symbols := parseSecurities(lines)
	return &SecuritySlice{symbols: symbols}
}

// NewCikCodeMap initializes a CikCodeMap
func NewCikCodeMap() *CikCodeMap {
	var ccm CikCodeMap
	ccm.urlEndpoint = urlString
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

func normalizeSymbol(symbol string) string {
	return strings.ReplaceAll(strings.Fields(symbol)[0], "/", ".")
}

func parseSecurities(lines [][]string) []string {
	var symbols []string
	exists := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		symb := normalizeSymbol(line[0])
		// skip expired
		if len(symb) >= 8 {
			continue
		}
		if _, found := exists[symb]; found {
			continue
		}
		symbols = append(symbols, symb)
		exists[symb] = struct{}{}
	}
	return symbols
}

// func main() {
// 	ccm := NewCikCodeMap()
// 	ccm.Build()
// 	aapl := ccm.Find("AAPL")
// 	fmt.Println(aapl)
//
// 	allsecs := NewSecuritySlice(csvfilename)
// 	fmt.Println(allsecs)
// }
