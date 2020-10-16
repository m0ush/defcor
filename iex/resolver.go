package iex

import "fmt"

func makeStockSets(A, B StockGroup) (map[Stock]struct{}, map[Stock]struct{}) {
	sA := make(map[Stock]struct{})
	for _, x := range A {
		sA[x] = struct{}{}
	}
	sB := make(map[Stock]struct{})
	for _, x := range B {
		sB[x] = struct{}{}
	}
	return sA, sB
}

func makeIntSet(xi []int) map[int]struct{} {
	si := make(map[int]struct{}, len(xi))
	for _, x := range xi {
		si[x] = struct{}{}
	}
	return si
}

func diffs(A, B StockGroup) (StockGroup, StockGroup) {
	sA, sB := makeStockSets(A, B)
	var ANotB StockGroup
	for _, x := range A {
		if _, found := sB[x]; !found {
			ANotB = append(ANotB, x)
		}
	}
	var BNotA StockGroup
	for _, x := range B {
		if _, found := sA[x]; !found {
			BNotA = append(BNotA, x)
		}
	}
	return ANotB, BNotA
}

func idMap(dA, dB StockGroup) (map[string]int, map[string]int) {
	mA := make(map[string]int, len(dA))
	for i, s := range dA {
		mA[s.IexID] = i
	}
	mB := make(map[string]int, len(dB))
	for i, s := range dB {
		mB[s.IexID] = i
	}
	return mA, mB
}

func findUpdates(mA, mB map[string]int) (map[int]int, []int, []int) {
	var ixA, ixB []int
	m := make(map[int]int)
	for k, va := range mA {
		if vb, ok := mB[k]; ok {
			m[va] = vb
			ixB = append(ixB, vb)
		} else {
			ixA = append(ixA, va)
		}
	}
	return m, ixA, ixB
}

func mapValues(m map[string]int) []int {
	var vs []int
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

func subset(b, bp []int) []int {
	sBp := makeIntSet(bp)
	var bnotbp []int
	for _, x := range b {
		if _, found := sBp[x]; !found {
			bnotbp = append(bnotbp, x)
		}
	}
	return bnotbp
}

func changes(m map[string]int, bp []int) []int {
	return subset(mapValues(m), bp)
}

func extractGroup(A StockGroup, xi []int) StockGroup {
	var corrected StockGroup
	for i := range xi {
		corrected = append(corrected, A[i])
	}
	return corrected
}

func extractFromMap(A, B StockGroup, t map[int]int) map[Stock]Stock {
	translator := make(map[Stock]Stock)
	for k, v := range t {
		translator[A[k]] = B[v]
	}
	return translator
}

// Resolve performs the Stock reconciliation
func Resolve(Existing, Refreshed StockGroup) (map[Stock]Stock, []Stock, []Stock) {
	dA, dB := diffs(Existing, Refreshed)
	mA, mB := idMap(dA, dB)
	translator, idxEnds, idxPotentialAdds := findUpdates(mA, mB)
	idxAdds := changes(mB, idxPotentialAdds)
	updates := extractFromMap(dA, dB, translator)
	deletes := extractGroup(dA, idxEnds)
	additions := extractGroup(dB, idxAdds)
	return updates, deletes, additions
}

// FormatOutput pretty prints the Resolve func
func FormatOutput(A, B StockGroup) {
	updates, deletes, additions := Resolve(A, B)
	fmt.Println("\nUpdates:")
	for k, v := range updates {
		fmt.Println("\nreconcile...")
		fmt.Println("\tprev:", k)
		fmt.Println("\tcurr:", v)
	}
	fmt.Println("\nEnds:")
	for i, s := range deletes {
		fmt.Printf("%2d: %v\n", i, s)
	}
	fmt.Println("\nAdds:")
	for i, s := range additions {
		fmt.Printf("%2d: %v\n", i, s)
	}
}
