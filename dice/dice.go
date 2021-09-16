package dice

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
)

var onlyOnce sync.Once

func Roll(number int, dice int, keep_highest int) int {

	var (
		seq     = make([]int, number)
		results []int
		result  int
	)

	zap.L().Debug("Rolling dices",
		zap.Int("Number of dices", number),
		zap.Int("Number of faces", dice),
		zap.Int("Keep Highest", keep_highest),
	)

	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano()) // only run once
	})

	for range seq {
		diceResult := rand.Intn(dice - 1)
		results = append(results, diceResult+1)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(results)))

	if keep_highest > 0 {
		result = Sum(results[:keep_highest])
	} else {
		result = Sum(results)
	}

	zap.L().Debug("Result of roll",
		zap.Int("result", result),
	)

	return result
}

func Sum(dices []int) int {
	var result = 0

	for _, value := range dices {
		result += value
	}

	return result
}

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
