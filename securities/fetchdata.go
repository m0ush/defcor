package fetcher

import (
	"encoding/csv"
	"os"
	"strings"
)

const csvfilename = "allsecurities.csv"

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
