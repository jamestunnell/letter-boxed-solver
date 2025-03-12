package models

import "github.com/yourbasic/bit"

type LetterSet struct {
	set *bit.Set
}

func NewLetterSet(words ...string) LetterSet {
	set := bit.New()

	for _, word := range words {
		for _, ch := range word {
			set.Add(int(ch))
		}
	}

	return LetterSet{set: set}
}

func (ls LetterSet) Size() int {
	return ls.set.Size()
}

func (ls LetterSet) EachRune(each func(r rune)) {
	ls.set.Visit(func(n int) (skip bool) {
		each(rune(n))

		return false
	})
}

func (ls LetterSet) String() string {
	runes := []rune{}

	ls.set.Visit(func(n int) (skip bool) {
		runes = append(runes, rune(n))

		return false
	})

	return string(runes)
}

func (ls LetterSet) Or(other LetterSet) LetterSet {
	return LetterSet{set: ls.set.Or(other.set)}
}

func (ls LetterSet) And(other LetterSet) LetterSet {
	return LetterSet{set: ls.set.And(other.set)}
}

func (ls LetterSet) AndNot(other LetterSet) LetterSet {
	return LetterSet{ls.set.AndNot(other.set)}
}
