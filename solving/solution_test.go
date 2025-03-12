package solving_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jamestunnell/letter-boxed-solver/solving"
)

func TestSolutionHash(t *testing.T) {
	s1 := solving.Solution{"APPLE", "EAGER", "RENT"}
	s2 := solving.Solution{"APPLE", "EAGER", "RENT"}

	assert.Equal(t, s1.Hash64(), s2.Hash64())
}
