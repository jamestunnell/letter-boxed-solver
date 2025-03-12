package solving

import (
	"slices"
	"sort"

	"github.com/jamestunnell/letter-boxed-solver/models"
	"github.com/jamestunnell/letter-boxed-solver/util"
)

type ExploreResults struct {
	Complete   []Solution
	Incomplete []Solution

	existing map[uint64]struct{}
}

type Explorer struct {
	*WordInfo

	puzzle      *models.Puzzle
	wordMapping *WordMapping
	maxBranch   int
	scoring     Scoring
}

func NewExploreResults() *ExploreResults {
	return &ExploreResults{
		Incomplete: []Solution{},
		Complete:   []Solution{},
		existing:   map[uint64]struct{}{},
	}
}

func NewExplorer(
	start *WordInfo,
	puzzle *models.Puzzle,
	wordMapping *WordMapping,
	maxBranch int,
	scoring Scoring,
) *Explorer {
	e := &Explorer{
		WordInfo:    start,
		puzzle:      puzzle,
		wordMapping: wordMapping,
		maxBranch:   maxBranch,
		scoring:     scoring,
	}

	return e
}

func (e *Explorer) Explore() []Solution {
	leftResults := e.exploreLeft()
	rightResults := e.exploreRight()
	complete := []Solution{}

	for _, backwardsSln := range leftResults.Complete {
		complete = append(complete, backwardsSln.Reverse())
	}

	complete = append(complete, rightResults.Complete...)

	for _, left := range leftResults.Incomplete {
		left = left.Reverse()

		for _, right := range rightResults.Incomplete {
			right = slices.Clone(right[1:])

			slns := e.findCrossingSolutions(left, right)
			if len(slns) > 0 {
				complete = append(complete, slns...)
			}
		}
	}

	return complete
}

func (e *Explorer) exploreLeft() *ExploreResults {
	results := NewExploreResults()

	e.explore([]*WordInfo{e.WordInfo}, e.Letters, e.getLeftSubwords, results)

	return results
}

func (e *Explorer) exploreRight() *ExploreResults {
	results := NewExploreResults()

	e.explore([]*WordInfo{e.WordInfo}, e.Letters, e.getRightSubwords, results)

	return results
}

func (e *Explorer) findCrossingSolutions(leftWords, rightWords []string) []Solution {
	solutions := []Solution{}

	for start := 0; start < (len(leftWords) - 1); start++ {
		words := slices.Clone(leftWords[start:])
		letters := models.NewLetterSet(words...)

		for _, rightWord := range rightWords {
			words = append(words, rightWord)
			letters = letters.Or(models.NewLetterSet(rightWord))

			if e.puzzle.DoLettersSolve(letters) {
				solutions = append(solutions, Solution(words))

				break
			} else if len(words) == e.puzzle.GetMaxWords() {
				break
			}
		}
	}

	return solutions
}

func (e *Explorer) getLeftSubwords(
	current *WordInfo,
	totalLetters models.LetterSet,
) []*WordInfo {
	subWords := e.wordMapping.WordsWithLastLetter(current.FirstLetter)

	if len(subWords) > e.maxBranch {
		subWords = e.reduceSubwords(subWords, totalLetters)
	}

	return subWords
}

func (e *Explorer) getRightSubwords(
	current *WordInfo,
	totalLetters models.LetterSet,
) []*WordInfo {
	subWords := e.wordMapping.WordsWithFirstLetter(current.LastLetter)

	if len(subWords) > e.maxBranch {
		subWords = e.reduceSubwords(subWords, totalLetters)
	}

	return subWords
}

func (e *Explorer) reduceSubwords(
	subWords []*WordInfo,
	totalLetters models.LetterSet,
) []*WordInfo {
	sortByScoreDesc := &SortWordsByScoreDesc{
		SortWordsByScore: &SortWordsByScore{
			Infos: subWords,
			Scores: util.Map(subWords, func(info *WordInfo) float64 {
				diff := info.Letters.AndNot(totalLetters)

				return e.scoring.Score(diff)
			}),
		},
	}

	sort.Sort(sortByScoreDesc)

	return subWords[:e.maxBranch]
}

func (e *Explorer) explore(
	current []*WordInfo,
	totalLetters models.LetterSet,
	getSubWords func(*WordInfo, models.LetterSet) []*WordInfo,
	results *ExploreResults,
) {
	if e.puzzle.DoLettersSolve(totalLetters) {
		results.AddComplete(util.Map(current, getWord))

		return
	} else if len(current) == (e.puzzle.GetMaxWords() - 1) {
		results.AddIncomplete(util.Map(current, getWord))
	}

	if len(current) == e.puzzle.GetMaxWords() {
		return
	}

	subWords := getSubWords(current[len(current)-1], totalLetters)
	for _, subWord := range subWords {
		// detect cycle
		if util.Any(current, func(info *WordInfo) bool {
			return info.Word == subWord.Word
		}) {
			continue
		}

		newTotalLetters := totalLetters.Or(subWord.Letters)

		e.explore(append(current, subWord), newTotalLetters, getSubWords, results)
	}
}

func getWord(info *WordInfo) string {
	return info.Word
}

func (results *ExploreResults) AddComplete(sln Solution) {
	hash := sln.Hash64()
	if _, found := results.existing[hash]; found {
		return
	}

	results.Complete = append(results.Complete, sln)
	results.existing[hash] = struct{}{}
}

func (results *ExploreResults) AddIncomplete(sln Solution) {
	hash := sln.Hash64()
	if _, found := results.existing[hash]; found {
		return
	}

	results.Incomplete = append(results.Incomplete, sln)
	results.existing[hash] = struct{}{}
}
