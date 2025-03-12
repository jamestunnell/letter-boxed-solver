package solving

import "slices"

type WordMapping struct {
	byFirstLetter map[rune][]*WordInfo
	byLastLetter  map[rune][]*WordInfo
}

func NewWordMapping(infos []*WordInfo) *WordMapping {
	byFirstLetter := map[rune][]*WordInfo{}
	byLastLetter := map[rune][]*WordInfo{}

	for _, info := range infos {
		byFirstLetter[info.FirstLetter] = append(byFirstLetter[info.FirstLetter], info)
		byLastLetter[info.LastLetter] = append(byLastLetter[info.LastLetter], info)
	}

	return &WordMapping{
		byFirstLetter: byFirstLetter,
		byLastLetter:  byLastLetter,
	}
}

func (wm *WordMapping) WordsWithFirstLetter(r rune) []*WordInfo {
	return slices.Clone(wm.byFirstLetter[r])
}

func (wm *WordMapping) WordsWithLastLetter(r rune) []*WordInfo {
	return slices.Clone(wm.byLastLetter[r])
}
