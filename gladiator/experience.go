package gladiator

import (
	"math"

	"github.com/wihrt/idle_arena/dice"
)

type Experience struct {
	Level     int `json:"level" bson:"level"`
	Current   int `json:"current" bson:"current"`
	NextLevel int `json:"next_level" bson:"next_level"`
}

func NewExperience() *Experience {
	e := &Experience{
		Level:     1,
		Current:   0,
		NextLevel: calculateNextLevel(1),
	}

	return e
}

func (g *Gladiator) LevelUp() {

	var HpBonus = dice.Roll(1, 12, -1)
	g.Health.Max += HpBonus + g.Constitution.Modifier
	g.Health.Current = g.Health.Max

	if (g.Experience.Level % 2) == 0 {
		stat := dice.Roll(1, 3, -1)
		switch stat {
		case 1:
			g.Strength.Add(1)
		case 2:
			g.Dexterity.Add(1)
			g.ArmorClass = calculateArmorClass(g)
		case 3:
			g.Constitution.Add(1)
		}
	}

	g.Experience.Level += 1
	remainingExp := g.Experience.Current - g.Experience.NextLevel
	g.Experience.NextLevel = calculateNextLevel(g.Experience.Level)
	g.Experience.Current = remainingExp
}

func calculateNextLevel(level int) int {
	var power = float64(1)
	if level > 1 {
		additional := float64(level-1) / 100
		power = float64(1) + additional
	}

	experienceNeeded := math.Floor(math.Pow(100, power))
	return int(experienceNeeded)
}
