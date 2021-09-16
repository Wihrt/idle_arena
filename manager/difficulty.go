package manager

import "errors"

type Difficulty int

const (
	DifficultyUnknown Difficulty = 0
	DifficultyEasy    Difficulty = 1
	DifficultyNormal  Difficulty = 2
	DifficultyHard    Difficulty = 3
)

func ParseDifficulty(difficulty int) (Difficulty, error) {

	switch difficulty {
	case 0:
		return DifficultyEasy, nil
	case 1:
		return DifficultyNormal, nil
	case 2:
		return DifficultyHard, nil
	default:
		return DifficultyUnknown, errors.New("unknown difficulty")
	}
}
