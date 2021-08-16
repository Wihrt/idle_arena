package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Difficulty int

const (
	DifficultyUnknown Difficulty = -1
	DifficultyEasy    Difficulty = 0
	DifficultyNormal  Difficulty = 1
	DifficultyHard    Difficulty = 2
)

type Manager struct {
	ManagerID  string     `json:"manager_id" bson:"manager_id"`
	Difficulty Difficulty `json:"difficulty" bson:"difficulty"`
	Gladiators []string   `json:"gladiators" bson:"gladiators"`
}

func NewManager(managerID string, difficulty int) (*Manager, error) {
	m := &Manager{
		ManagerID: managerID,
	}

	d, err := ParseDifficulty(difficulty)
	if err != nil {
		return m, err
	}

	m.Difficulty = d

	return m, nil
}

func GenerateID(guildID string, userID string) string {
	h := sha256.Sum256([]byte(guildID + "." + userID))
	id := hex.EncodeToString(h[:])
	return id
}

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
