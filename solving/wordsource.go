package solving

import (
	"bufio"
	"io/fs"
)

type WordSource interface {
	NextWord() (string, bool)
}

type FileWordSource struct {
	scanner *bufio.Scanner
}

func NewFileWordSource(f fs.File) WordSource {
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	return &FileWordSource{
		scanner: scanner,
	}
}

func (ws *FileWordSource) NextWord() (string, bool) {
	if !ws.scanner.Scan() {
		return "", false
	}

	return ws.scanner.Text(), true
}
