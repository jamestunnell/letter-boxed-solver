package solving

import "github.com/jamestunnell/letter-boxed-solver/models"

type Scoring interface {
	Score(models.LetterSet) float64
}

type UniformScoring struct {
}

type WeightedScoring struct {
	LetterWeights map[rune]float64
}

func NewWeightedScoring(infos []*WordInfo) *WeightedScoring {
	inWordsTotals := map[rune]int{}

	for _, info := range infos {
		info.Letters.EachRune(func(r rune) {
			inWordsTotals[r]++
		})
	}

	numWords := float64(len(infos))
	weights := map[rune]float64{}
	for r, total := range inWordsTotals {
		weights[r] = float64(total) / numWords
	}

	return &WeightedScoring{
		LetterWeights: weights,
	}
}

func (s *UniformScoring) Score(ls models.LetterSet) float64 {
	return float64(ls.Size())
}

func (s *WeightedScoring) Score(ls models.LetterSet) float64 {
	score := 0.0
	ls.EachRune(func(r rune) {
		if weight, found := s.LetterWeights[r]; found {
			score += weight
		}
	})

	return float64(ls.Size())
}
