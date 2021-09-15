package fight

import (
	"math"

	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/manager"
)

func ExperienceGained(m *manager.Manager, s *Settings) int {

	var (
		managerMultiplier float64
		experienceRoll    int
		experienceGained  float64
	)

	switch m.Difficulty {
	case manager.DifficultyEasy:
		managerMultiplier = 0.5
	case manager.DifficultyNormal:
		managerMultiplier = 1
	case manager.DifficultyHard:
		managerMultiplier = 1.5
	}

	experienceRoll = dice.Roll(int(s.Difficulty)+1, 20, -1)
	experienceGained = math.Floor(float64(experienceRoll) * managerMultiplier)

	return int(experienceGained)
}

func GoldGained(m *manager.Manager, s *Settings) int {

	var (
		fightMultiplier   float64
		managerMultiplier float64
		goldRoll          int
		goldGained        float64
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

	switch m.Difficulty {
	case manager.DifficultyEasy:
		managerMultiplier = 0.5
	case manager.DifficultyNormal:
		managerMultiplier = 1
	case manager.DifficultyHard:
		managerMultiplier = 1.5
	}

	goldRoll = dice.Roll(1, 10, -1)
	goldGained = math.Floor(float64(goldRoll) * fightMultiplier * managerMultiplier)
	return int(goldGained)
}
