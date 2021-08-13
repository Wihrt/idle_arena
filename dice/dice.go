package dice

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

var onlyOnce sync.Once

func Roll(number int, dice int, keep_highest int) int {

	var (
		seq     = make([]int, number)
		results []int
		result  int
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
