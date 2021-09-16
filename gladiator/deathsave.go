package gladiator

import "github.com/wihrt/idle_arena/manager"

type DeathSave struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

func NewDeathSave(difficulty manager.Difficulty) *DeathSave {
	d := &DeathSave{
		Current: 0,
	}

	switch difficulty {
	case manager.DifficultyEasy:
		d.Max = 6
	case manager.DifficultyNormal:
		d.Max = 4
	case manager.DifficultyHard:
		d.Max = 2
	}
	return d
}
