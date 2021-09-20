package manager

import "errors"

type Difficulty int

const (
	DifficultyUnknown Difficulty = 0
	DifficultyEasy    Difficulty = 1
	DifficultyNormal  Difficulty = 2
	DifficultyHard    Difficulty = 3
)

func ParseDifficulty(difficulty int) (Difficulty, float64, error) {

	switch difficulty {
	case 0:
		return DifficultyEasy, 0.5, nil
	case 1:
		return DifficultyNormal, 1, nil
	case 2:
		return DifficultyHard, 1.5, nil
	default:
		return DifficultyUnknown, 1, errors.New("unknown difficulty")
	}
}
