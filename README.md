# letter-boxed-solver
Command-line tool to solve puzzles like the NYT letter boxed puzzle. The tool has a few commands: `list-builtin`, `solve-builtin`, and `solve-given`.

Built-in puzzles filenames can be listed with the `list-builtin` command. 

Then a built-in puzzle can be solved using `solve-builtin`:
```console
Usage: letter-boxed-solver solve-builtin [--maxbranch MAXBRANCH] [--maxtime MAXTIME] [--outdir OUTDIR] --fname FNAME

Options:
  --maxbranch MAXBRANCH
                         max degree of a solving branch [default: 5]
  --maxtime MAXTIME      max time to spend solving [default: 5s]
  --outdir OUTDIR, -o OUTDIR
                         output directory (created if it does not exist) [default: .]
  --fname FNAME          name of a built-in puzzle file
```

Alternately, a puzzle can be given via the command line with `solve-given`:
```console
Usage: letter-boxed-solver solve-given [--maxbranch MAXBRANCH] [--maxtime MAXTIME] [--outdir OUTDIR] --maxwords MAXWORDS --sides SIDES

Options:
  --maxbranch MAXBRANCH
                         max degree of a solving branch [default: 5]
  --maxtime MAXTIME      max time to spend solving [default: 5s]
  --outdir OUTDIR, -o OUTDIR
                         output directory (created if it does not exist) [default: .]
  --maxwords MAXWORDS    Max words to allow for puzzle solution
  --sides SIDES          Puzzle sides, with letters combined (e.g. abc def ghi jkl)
```