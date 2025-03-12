package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/yourbasic/bit"
)

type Puzzle struct {
	sides           []string
	maxWords        int
	letterSet       *bit.Set
	antiConnections map[rune][]rune
}

type PuzzleData struct {
	MaxWords int      `json:"maxWords"`
	Sides    []string `json:"sides"`
}

type Side []rune

var errLetterRepeated = errors.New("letter is repeated")

func (p *Puzzle) UnmarshalJSON(d []byte) error {
	var pd PuzzleData

	if err := json.Unmarshal(d, &pd); err != nil {
		return err
	}

	letters := map[rune]struct{}{}
	for _, side := range pd.Sides {
		for _, letter := range side {
			if _, found := letters[letter]; found {
				return fmt.Errorf("invalid puzzle: %w", errLetterRepeated)
			}

			letters[letter] = struct{}{}
		}
	}

	p.Init(&pd)

	return nil
}

func (p *Puzzle) Init(pd *PuzzleData) {
	antiConnections := map[rune][]rune{}
	letterSet := bit.New()

	for _, side := range pd.Sides {
		side = strings.ToUpper(side)

		for _, ch := range side {
			letterSet.Add(int(ch))

			antiConnections[ch] = []rune(side)
		}
	}

	p.sides = pd.Sides
	p.antiConnections = antiConnections
	p.letterSet = letterSet
	p.maxWords = pd.MaxWords
}

func (p *Puzzle) GetMaxWords() int {
	return p.maxWords
}

func (p *Puzzle) IsWordAllowed(word string) bool {
	const minWordLen = 3

	if len(word) < minWordLen {
		return false
	}

	disallowed := []rune{}

	for _, ch := range word {
		if slices.Contains(disallowed, ch) {
			return false
		}

		var found bool

		if disallowed, found = p.antiConnections[ch]; !found {
			return false
		}
	}

	return true
}

func (p *Puzzle) GetSides() []string {
	return p.sides
}

func (p *Puzzle) GetLetterSet() LetterSet {
	return LetterSet{set: p.letterSet}
}

func (p *Puzzle) DoLettersSolve(ls LetterSet) bool {
	return ls.set.And(p.letterSet).Size() == p.letterSet.Size()
}
