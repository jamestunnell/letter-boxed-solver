package solving

import (
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/letter-boxed-solver/models"
	"github.com/jamestunnell/letter-boxed-solver/util"
)

type Solver struct {
	puzzle    *models.Puzzle
	explorers []*Explorer
	scoring   Scoring
	solutions []Solution
	existing  map[uint64]struct{}
}

func NewSolver(
	p *models.Puzzle,
	wordSource WordSource,
	maxBranch int,
) *Solver {
	allowedWords := loadWords(wordSource, p.IsWordAllowed)

	log.Info().Msg("making word graph")

	solutions := []Solution{}
	unsolved := []*WordInfo{}
	for _, word := range allowedWords {
		info := NewWordInfo(word)
		if p.DoLettersSolve(info.Letters) {
			solutions = append(solutions, Solution{word})
		} else {
			unsolved = append(unsolved, info)
		}
	}

	wm := NewWordMapping(unsolved)
	scoring := NewWeightedScoring(unsolved)
	sortByScoreAsc := &SortWordsByScoreAsc{
		SortWordsByScore: &SortWordsByScore{
			Infos: unsolved,
			Scores: util.Map(unsolved, func(info *WordInfo) float64 {
				return scoring.Score(info.Letters)
			}),
		},
	}

	// sort so the best prospect is at the end
	sort.Sort(sortByScoreAsc)

	return &Solver{
		puzzle: p,
		explorers: util.Map(unsolved, func(info *WordInfo) *Explorer {
			return NewExplorer(info, p, wm, maxBranch, scoring)
		}),
		scoring:   scoring,
		solutions: solutions,
		existing:  map[uint64]struct{}{},
	}
}

func (s *Solver) IsFinished() bool {
	return len(s.explorers) == 0
}

func (s *Solver) GetSolutions() []Solution {
	return s.solutions
}

func (s *Solver) Step() {
	remaining := len(s.explorers)

	if remaining == 0 {
		return
	}

	e := s.explorers[remaining-1]

	solutions := e.Explore()
	for _, sln := range solutions {
		hash := sln.Hash64()
		if _, found := s.existing[hash]; found {
			continue
		}

		s.solutions = append(s.solutions, sln)
		s.existing[hash] = struct{}{}
	}

	// pop
	s.explorers = s.explorers[:remaining-1]
}

func loadWords(source WordSource, allowed func(string) bool) []string {
	total := 0
	words := []string{}

	nextWord, ok := source.NextWord()
	for ok {
		if allowed(nextWord) {
			words = append(words, nextWord)
		}

		total++

		nextWord, ok = source.NextWord()
	}

	log.Info().
		Int("total", total).
		Int("valid", len(words)).
		Msg("loaded words file")

	return words
}
