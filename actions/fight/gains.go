package fight

import (
	"math"

	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/manager"
)

func fightMultiplier(s *Settings) float64 {

	var fightMultiplier float64

	switch s.Difficulty {
	case DifficultyEasy:
		fightMultiplier = 0.6
	case DifficultyNormal:
		fightMultiplier = 0.8
	case DifficultyHard:
		fightMultiplier = 1
	case DifficultyChallenging:
		fightMultiplier = 1.2
	case DifficultyHellish:
		fightMultiplier = 1.4
	case DifficultyNightmarish:
		fightMultiplier = 1.6
	}

	return fightMultiplier
}

func ExperienceGained(m *manager.Manager, s *Settings) int {

	var (
		experienceRoll   int
		experienceGained float64
	)

	experienceRoll = dice.Roll(int(s.Difficulty)+1, 20, -1)
	experienceGained = math.Floor(float64(experienceRoll) * m.Multiplier)

	return int(experienceGained)
}

func GoldGained(m *manager.Manager, s *Settings) int {

	var (
		fightMultiplier = fightMultiplier(s)
		goldRoll        int
		goldGained      float64
	)

	switch s.Difficulty {
	case DifficultyEasy:
		fightMultiplier = 0.6
	case DifficultyNormal:
		fightMultiplier = 0.8
	case DifficultyHard:
		fightMultiplier = 1
	case DifficultyChallenging:
		fightMultiplier = 1.2
	case DifficultyHellish:
		fightMultiplier = 1.4
	case DifficultyNightmarish:
		fightMultiplier = 1.6
	}

	goldRoll = dice.Roll(1, 10, -1)
	goldGained = math.Floor(float64(goldRoll) * fightMultiplier * m.Multiplier)
	return int(goldGained)
}
