package gladiator

import (
	"math"
)

func calculateNextLevel(level int) int {
	var power = float64(1)
	if level > 1 {
		additional := float64(level-1) / 100
		power = float64(1) + additional
	}

	experienceNeeded := math.Floor(math.Pow(100, power))
	return int(experienceNeeded)
}
