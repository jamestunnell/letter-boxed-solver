package solving

type SortWordsByScore struct {
	Infos  []*WordInfo
	Scores []float64
}

type SortWordsByScoreAsc struct {
	*SortWordsByScore
}

type SortWordsByScoreDesc struct {
	*SortWordsByScore
}

func (s *SortWordsByScore) Len() int {
	return len(s.Infos)
}

func (s *SortWordsByScoreAsc) Less(i, j int) bool {
	return s.Scores[i] < s.Scores[j]
}

func (s *SortWordsByScoreDesc) Less(i, j int) bool {
	return s.Scores[i] > s.Scores[j]
}

func (s *SortWordsByScore) Swap(i, j int) {
	s.Scores[i], s.Scores[j] = s.Scores[j], s.Scores[i]
	s.Infos[i], s.Infos[j] = s.Infos[j], s.Infos[i]
}
