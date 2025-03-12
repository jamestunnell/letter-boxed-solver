package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/letter-boxed-solver/models"
	"github.com/jamestunnell/letter-boxed-solver/solving"
)

//go:embed puzzles/*
var puzzles embed.FS

//go:embed words/*
var words embed.FS

type ListBuiltinCmd struct {
}

type SolveBuiltinCmd struct {
	MaxBranch int    `help:"max degree of a solving branch" default:"5"`
	MaxTime   string `help:"max time to spend solving" default:"250ms"`
	Filename  string `arg:"-f,required" help:"name of a built-in puzzle file"`
	Outdir    string `arg:"-o" help:"output directory (created if it does not exist)" default:"."`
}

type Args struct {
	ListBuiltin  *ListBuiltinCmd  `arg:"subcommand:list-builtin" help:"list built-in puzzle files"`
	SolveBuiltin *SolveBuiltinCmd `arg:"subcommand:solve-builtin" help:"solve built-in puzzle file"`
}

func (args Args) Version() string {
	return "0.1.0"
}

func main() {
	var args Args

	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		err = fmt.Errorf("failed to create arg parser: %w", err)

		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	err = p.Parse(os.Args[1:])

	switch {
	case errors.Is(err, arg.ErrVersion): // found "--version" on command line
		fmt.Println(args.Version())
		os.Exit(0)
	case errors.Is(err, arg.ErrHelp): // found "--help" on command line
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	case err != nil:
		fmt.Printf("error: %v\n\n", err)
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if p.Subcommand() == nil {
		fmt.Printf("error: missing command\n\n")
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	switch {
	case args.ListBuiltin != nil:
		err = listBuiltin()
	case args.SolveBuiltin != nil:
		err = solveBuiltin(args.SolveBuiltin)
	default:
	}

	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}
}

func listBuiltin() error {
	entries, err := puzzles.ReadDir("puzzles")
	if err != nil {
		return fmt.Errorf("failed to failed to read puzzle entries: %w", err)
	}

	log.Info().Int("count", len(entries)).Msg("found puzzle entries")

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	return nil
}

func solveBuiltin(cmd *SolveBuiltinCmd) error {
	fname := fmt.Sprintf("puzzles/%s", cmd.Filename)
	name := strings.TrimSuffix(cmd.Filename, path.Ext(cmd.Filename))
	outpath := fmt.Sprintf("%s/%s-solutions.txt", cmd.Outdir, name)

	maxTime, err := time.ParseDuration(cmd.MaxTime)
	if err != nil {
		return fmt.Errorf("error: failed to parse max time: %w", err)
	}

	f, err := puzzles.Open(fname)
	if err != nil {
		return fmt.Errorf("error: failed to open built-in puzzle file: %w", err)
	}

	log.Info().
		Str("maxTime", cmd.MaxTime).
		Str("fname", fname).
		Str("outpath", outpath).
		Msg("solving built-in puzzle")

	solutions, err := solve(f, maxTime, cmd.MaxBranch)
	if err != nil {
		return err
	}

	if err = reportSolutions(solutions, outpath); err != nil {
		return err
	}

	return nil
}

func loadPuzzle(puzzleFile fs.File) (*models.Puzzle, error) {
	var p models.Puzzle

	if err := json.NewDecoder(puzzleFile).Decode(&p); err != nil {
		return nil, fmt.Errorf("failed to load puzzle JSON file: %w", err)
	}

	return &p, nil
}

func solve(
	puzzleFile fs.File,
	maxTime time.Duration,
	maxBranch int,
) (solving.SolutionsByWordCount, error) {
	puzzle, err := loadPuzzle(puzzleFile)
	if err != nil {
		return solving.SolutionsByWordCount{}, err
	}

	log.Info().
		Strs("sides", puzzle.GetSides()).
		Stringer("letters", puzzle.GetLetterSet()).
		Int("maxWords", puzzle.GetMaxWords()).
		Msg("loaded puzzle")

	wordsFile, err := words.Open("words/scrabble-words.txt")
	if err != nil {
		return solving.SolutionsByWordCount{}, fmt.Errorf("failed to open words file: %w", err)
	}

	wordSource := solving.NewFileWordSource(wordsFile)

	start := time.Now()

	log.Info().
		Float64("maxTimeSec", maxTime.Seconds()).
		Msg("solving puzzle")

	solver := solving.NewSolver(puzzle, wordSource, maxBranch)
	solutions := solving.SolutionsByWordCount{}
	step := 0

	for !solver.IsFinished() && (time.Since(start) <= maxTime) {
		solver.Step()

		fmt.Print(".")

		step++
	}

	log.Info().
		Float64("durSec", time.Since(start).Seconds()).
		Msg("done solving")

	for _, sln := range solver.GetSolutions() {
		solutions.Add(sln)
	}

	return solutions, nil
}

func reportSolutions(
	solutions solving.SolutionsByWordCount,
	outpath string,
) error {
	allSlns := solutions.All()

	if len(allSlns) > 0 {
		log.Info().
			Int("count", len(allSlns)).
			Strs("best", allSlns[0]).
			Msg("found solutions")
	} else {
		log.Info().Msg("no solutions found")
	}

	outdir := path.Dir(outpath)

	info, err := os.Stat(outdir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat output dir: %w", err)
		}

		if err = os.Mkdir(outdir, 0750); err != nil {
			return fmt.Errorf("failed to make output dir: %w", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("'%s' is not a dir", outdir)
	}

	solutionsFile, err := os.Create(outpath)
	if err != nil {
		return fmt.Errorf("failed to create solutions file: %w", err)
	}

	log.Info().Str("outpath", outpath).Msg("writing solutions to file")

	defer solutionsFile.Close()

	w := bufio.NewWriter(solutionsFile)

	defer w.Flush()

	for _, sln := range allSlns {
		w.WriteString(sln.String())
		w.WriteRune('\n')
	}

	return nil
}
