package skiplist

import (
	"math/rand"
	"sync"
	"time"
)

const (
	p = .5 // the p level defines the probability that a node
	// with a value at level i also has a value at i+1.  This number
	// is also important in determining max level.  Max level will
	// be defined as L(N) where L = log base (1/p) of n where n
	// is the number of items in the list and N is the number of possible
	// items in the universe.  If p = .5 then maxlevel = 32 is appropriate
	// for uint32.
	maxlevel = 32
)

// lockedSource is an implementation of rand.Source that is safe for
// concurrent use by multiple goroutines. The code is modeled after
// https://golang.org/src/math/rand/rand.go.
type lockedSource struct {
	mu  sync.Mutex
	src rand.Source
}

// Int63 implements the rand.Source interface.
func (ls *lockedSource) Int63() (n int64) {
	ls.mu.Lock()
	n = ls.src.Int63()
	ls.mu.Unlock()
	return
}

// Seed implements the rand.Source interface.
func (ls *lockedSource) Seed(seed int64) {
	ls.mu.Lock()
	ls.src.Seed(seed)
	ls.mu.Unlock()
}

// generator will be the common generator to create random numbers. It
// is seeded with unix nanosecond when this line is executed at runtime,
// and only executed once ensuring all random numbers come from the same
// randomly seeded generator.
var generator = rand.New(&lockedSource{src: rand.NewSource(time.Now().UnixNano())})

func flipCoin() bool {
	return generator.Float64() >= p
}

func generateLevel(maxLevel int) (level int) {
	for level = 1; level < maxLevel-1 && flipCoin(); level++ {
	}
	return level
}
