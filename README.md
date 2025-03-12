# letter-boxed-solver
Command-line tool to solve puzzles like the NYT letter boxed puzzle.

## Puzzle Rules
The NYT letter-boxed puzzle has four sides with three letters, forming a box. The letters are not repeated. The puzzle is solved using a list of words. For each word, there must not be two consecutive letters from the same side. The first letter of each word after the first word must match the last letter of the previous word. Also, a puzzle has a max solution word count.

## Solver

The solver attempts to find all possible solution, starting with those that should lead to the shortest solutions.

The solver begins by loading a dictionary of words that *might* be accepted by the actual NYT puzzle app/website. Then the puzzle is used to determine which of these words are allowed.

The allowed words are sorted by score. Searching proceeds in order from highest to lowest score. Words are scored according to how many and which puzzle letters they will visit. A weighted scoring system favors letters that occur in fewer of the allowed words. The total word score is the weighted sum of letter scores.

Potential solutions are explored by branching to the left and right of a start word. Branching can be limited with a parameter, and only the highest value sub-word paths are explored. A sub-word is scored only using the letters that do not already exist in any words in the current exploration path. Exploration along a path will end early if a complete solution is found. Exploration also ends when max solution words are reached without having a complete solution.

Once the solutions around an allowed word have been explored, there will be a number of incomplete word chains to the left and right. Some of these can be used to form complete solutions by crossing left-to-right past the start word.

## Usage

Built-in puzzles filenames can be listed with the `list-builtin` command. Then a built-in puzzle can be solved using `solve-builtin`. Alternately, a puzzle can be given via the command line with `solve-given`.

When solving, the `--maxtime` option can be used to limit solutions to those that would be found soonest. Since the solver starts with highest-value words, this will most likely include the overall best solution. The `--maxbranch` option will limit the amount of branching as potential solutions are explored. Since only the highest-value sub-words would be used for exploration, a smaller max branch value will probably not prevent reaching the overall best solution.

Solutions are written to the output file, one solution per line, sorted by number of words (ascending) and total characters (ascending).
