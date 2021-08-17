package fight

import "github.com/wihrt/idle_arena/manager"

func ExperienceGained(m *manager.Manager, s *Settings) int {

	switch m.Difficulty {
	case manager.DifficultyEasy:
	case manager.DifficultyNormal:
	case manager.DifficultyHard:

	}
}

func GoldGained(m *manager.Manager, s *Settings) int {

	switch s.Difficulty {
	case DifficultyEasy:
	case DifficultyNormal:
	case DifficultyHard:
	case DifficultyChallenging:
	case DifficultyHellish:
	case DifficultyNightmarish:
	}

	switch m.Difficulty {
	case manager.DifficultyEasy:
	case manager.DifficultyNormal:
	case manager.DifficultyHard:

	}

}
