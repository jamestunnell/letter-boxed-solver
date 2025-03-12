package solving

import (
	"cmp"
	"hash/maphash"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type Solution []string
type SolutionsByWordCount map[int][]Solution

var seed maphash.Seed

func init() {
	seed = maphash.MakeSeed()
}

// type ByWordCount Solutions

// func (wc ByWordCount) Len() int      { return len(wc) }
// func (wc ByWordCount) Swap(i, j int) { wc[i], wc[j] = wc[j], wc[i] }
// func (wc ByWordCount) Less(i, j int) bool {
// 	return cmp.Less(wc[i].WordCount(), wc[j].WordCount())
// }

// type ByTotalChars Solutions

// func (tc ByTotalChars) Len() int      { return len(tc) }
// func (tc ByTotalChars) Swap(i, j int) { tc[i], tc[j] = tc[j], tc[i] }
// func (tc ByTotalChars) Less(i, j int) bool {
// 	return cmp.Less(tc[i].TotalChars(), tc[j].TotalChars())
// }

func (slnsByWC SolutionsByWordCount) All() []Solution {
	all := []Solution{}
	wordCounts := maps.Keys(slnsByWC)

	slices.Sort(wordCounts)

	for _, wc := range wordCounts {
		slns := slnsByWC[wc]

		slices.SortFunc(slns, func(a, b Solution) int {
			return cmp.Compare(a.TotalChars(), b.TotalChars())
		})

		all = append(all, slns...)
	}

	return all
}

func (slnsByWC SolutionsByWordCount) Add(newSln Solution) {
	wc := len(newSln)

	slns, found := slnsByWC[wc]
	if !found {
		slnsByWC[wc] = []Solution{newSln}

		return
	}

	slnsByWC[wc] = append(slns, newSln)
}

func (s Solution) Hash64() uint64 {
	var hash maphash.Hash

	hash.SetSeed(seed)

	for _, word := range s {
		hash.WriteString(word)
	}

	return hash.Sum64()
}

func (s Solution) TotalChars() int {
	totalChars := 0

	for _, word := range s {
		totalChars += len(word)
	}

	return totalChars
}

func (s Solution) String() string {
	return strings.Join(s, ", ")
}

func (s Solution) Reverse() Solution {
	reversed := slices.Clone(s)

	slices.Reverse(reversed)

	return reversed
}
