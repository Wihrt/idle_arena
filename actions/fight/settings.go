package fight

import (
	"errors"

	"go.uber.org/zap"
)

type Difficulty int

const (
	DifficultyEasy        = 0
	DifficultyNormal      = 1
	DifficultyHard        = 2
	DifficultyChallenging = 3
	DifficultyNightmarish = 4
	DifficultyHellish     = 5
)

type Settings struct {
	Difficulty Difficulty `json:"difficulty"`
}

func NewSettings(difficulty string) (*Settings, error) {

	var s = &Settings{}

	d, err := ParseDifficulty(difficulty)
	if err != nil {
		zap.L().Error("Unknown difficulty",
			zap.String("difficulty", difficulty),
			zap.Error(err))
		return s, err
	}

	s.Difficulty = d

	return s, nil
}

func ParseDifficulty(difficulty string) (Difficulty, error) {

	var d Difficulty

	switch difficulty {
	case "easy":
		return DifficultyEasy, nil
	case "normal":
		return DifficultyNormal, nil
	case "hard":
		return DifficultyHard, nil
	case "challenging":
		return DifficultyChallenging, nil
	case "nightmarish":
		return DifficultyNightmarish, nil
	case "hellish":
		return DifficultyHellish, nil
	default:
		return d, errors.New("unknown difficulty")
	}
}
