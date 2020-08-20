package securities

import (
	"encoding/csv"
	"os"
	"strings"
)

const csvfilename = "securities/allsecurities.csv"

// Securities is a slice of security symbols
type Securities []string

// NewSecuritySlice instatiates the security slice
func NewSecuritySlice() (Securities, error) {
	file, err := os.Open(csvfilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	return parseSecurities(lines), nil
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
