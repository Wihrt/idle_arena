package gladiator

import (
	"math"

	"github.com/wihrt/idle_arena/dice"
)

func (g *Gladiator) LevelUp() {

	var HpBonus = dice.Roll(1, 12, -1)
	g.MaxHealth += HpBonus + g.Constitution.Modifier
	g.CurrentHealth = g.MaxHealth

	if (g.Level % 2) == 0 {
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

	g.Level += 1
	remainingExp := g.Experience - g.ExperienceToNextLevel
	g.ExperienceToNextLevel = calculateNextLevel(g.Level)
	g.Experience = remainingExp
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
