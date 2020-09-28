// Package weightedrand contains a performant data structure and algorithm used
// to randomly select an element from some kind of list, where the chances of
// each element to be selected not being equal, but defined by relative
// "weights" (or probabilities). This is called weighted random selection.
//
// Compare this package with (github.com/jmcvetta/randutil).WeightedChoice,
// which is optimized for the single operation case. In contrast, this package
// creates a presorted cache optimized for binary search, allowing for repeated
// selections from the same set to be significantly faster, especially for large
// data sets.
package weightedrand

import (
	"math/rand"
	"sort"
)

// Choice is a generic wrapper that can be used to add weights for any item.
type Choice struct {
	Item   interface{}
	Weight uint
}

// NewChoice creates a new Choice with specified item and weight.
func NewChoice(item interface{}, weight uint) Choice {
	return Choice{Item: item, Weight: weight}
}

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser struct {
	data   []Choice
	totals []int
	max    int
}

// NewChooser initializes a new Chooser for picking from the provided Choices.
func NewChooser(cs ...Choice) Chooser {
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Weight < cs[j].Weight
	})
	totals := make([]int, len(cs))
	runningTotal := 0
	for i, c := range cs {
		runningTotal += int(c.Weight)
		totals[i] = runningTotal
	}
	return Chooser{data: cs, totals: totals, max: runningTotal}
}

// Pick returns a single weighted random Choice.Item from the Chooser.
//
// Utilizes global rand as the source of randomness -- you will likely want to
// seed it.
func (chs Chooser) Pick() interface{} {
	r := rand.Intn(chs.max) + 1
	i := sort.SearchInts(chs.totals, r)
	return chs.data[i].Item
}

// PickSource returns a single weighted random Choice.Item from the Chooser,
// utilizing the provided *rand.Rand source rs for randomness.
//
// The primary use-case for this is avoid lock contention from the global random
// source if utilizing Chooser(s) from multiple goroutines in extremely
// high-throughput situations.
//
// It is the responsibility of the caller to ensure the provided rand.Source is
// safe from thread safety issues.
func (chs Chooser) PickSource(rs *rand.Rand) interface{} {
	r := rs.Intn(chs.max) + 1
	i := sort.SearchInts(chs.totals, r)
	return chs.data[i].Item
}
