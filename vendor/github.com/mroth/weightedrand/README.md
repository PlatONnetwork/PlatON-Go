# weightedrand :balance_scale:

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mroth/weightedrand)](https://pkg.go.dev/github.com/mroth/weightedrand)
[![CodeFactor](https://www.codefactor.io/repository/github/mroth/weightedrand/badge)](https://www.codefactor.io/repository/github/mroth/weightedrand)
[![Build Status](https://github.com/mroth/weightedrand/workflows/test/badge.svg)](https://github.com/mroth/weightedrand/actions)
[![codecov](https://codecov.io/gh/mroth/weightedrand/branch/master/graph/badge.svg)](https://codecov.io/gh/mroth/weightedrand)

> Fast weighted random selection for Go.

Randomly selects an element from some kind of list, where the chances of each
element to be selected are not equal, but rather defined by relative "weights"
(or probabilities). This is called weighted random selection.

## Usage

```go
import (
    /* ...snip... */
    wr "github.com/mroth/weightedrand"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

    c := wr.NewChooser(
        wr.Choice{Item: "🍒", Weight: 0},
        wr.Choice{Item: "🍋", Weight: 1},
        wr.Choice{Item: "🍊", Weight: 1},
        wr.Choice{Item: "🍉", Weight: 3},
        wr.Choice{Item: "🥑", Weight: 5},
    )
    /* The following will print 🍋 and 🍊 with 0.1 probability, 🍉 with 0.3
    probability, and 🥑 with 0.5 probability. 🍒 will never be printed. (Note
    the weights don't have to add up to 10, that was just done here to make the
    example easier to read.) */
    result := c.Pick().(string)
    fmt.Println(result)
}
```

## Benchmarks

The existing Go library that has a comparable implementation of this is
[`github.com/jmcvetta/randutil`][1], which optimizes for the single operation
case. In contrast, this library creates a presorted cache optimized for binary
search, allowing repeated selections from the same set to be significantly
faster, especially for large data sets.

[1]: https://github.com/jmcvetta/randutil

Comparison of this library versus `randutil.ChooseWeighted`. For repeated
samplings from large collections, `weightedrand` will be much quicker.

| Num choices |    `randutil` | `weightedrand` |
| ----------: | ------------: | -------------: |
|          10 |     435 ns/op |       58 ns/op |
|         100 |     511 ns/op |       84 ns/op |
|       1,000 |    1297 ns/op |      112 ns/op |
|      10,000 |    7952 ns/op |      137 ns/op |
|     100,000 |   85142 ns/op |      173 ns/op |
|   1,000,000 | 2082248 ns/op |      312 ns/op |

Don't be mislead by these numbers into thinking `weightedrand` is always the
right choice! If you are only picking from the same distribution once,
`randutil` will be faster. `weightedrand` optimizes for repeated calls at the
expense of some setup time and memory storage.

*Update: Starting in `v0.3.0` weightedrand can now scale linearly to take
advantage of multiple CPU cores in parallel, making it even faster. See
[PR#2](https://github.com/mroth/weightedrand/pull/2) for details.*

## Caveats

Note this uses `math/rand` instead of `crypto/rand`, as it is optimized for
performance, not a cryptographically secure implementation.

## Credits

To better understand the algorithm used in this library (as well as the one used
in randutil) check out this great blog post: [Weighted random generation in Python](https://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/).
