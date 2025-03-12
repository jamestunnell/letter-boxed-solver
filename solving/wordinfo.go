package solving

import (
	"github.com/jamestunnell/letter-boxed-solver/models"
	"github.com/rs/zerolog/log"
)

// type WordInfo interface {
// 	GetWord()                    string
// 	GetFirstLetter() rune
// 	GetLastLetter() rune
// 	GetLetterSet()               models.LetterSet
// 	GetTotalLetterSet() models.LetterSet
// }

type WordInfo struct {
	Word                    string
	FirstLetter, LastLetter rune
	Letters                 models.LetterSet
}

func NewWordInfo(word string) *WordInfo {
	runes := []rune(word)
	n := len(runes)

	if n == 0 {
		log.Fatal().Msg("unexpected empty word")
	}

	// var hash maphash.Hash

	// hash.WriteString(word)

	return &WordInfo{
		Word:        word,
		FirstLetter: runes[0],
		LastLetter:  runes[n-1],
		Letters:     models.NewLetterSet(word),
	}
}
