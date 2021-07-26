package gladiator

import "github.com/wihrt/idle_arena/dice"

type Gladiator struct {
	Armor                 *Armor         `json:"armor"`
	ArmorClass            int            `json:"armor_class"`
	Constitution          *Caracteristic `json:"constitution"`
	CurrentHealth         int            `json:"current_health"`
	Dexterity             *Caracteristic `json:"dexterity"`
	Experience            int            `json:"experience"`
	ExperienceToNextLevel int            `json:"experience_to_next_level"`
	Level                 int            `json:"level"`
	MaxHealth             int            `json:"max_health"`
	Name                  string         `json:"name"`
	Strength              *Caracteristic `json:"strength"`
	Weapon                *Weapon        `json:"weapon"`
}

func NewGladiator(level int) *Gladiator {
	g := &Gladiator{
		Level:                 1,
		Experience:            0,
		ExperienceToNextLevel: calculateNextLevel(1),
	}

	g.Strength = NewCaracteristic("strength", 4, 6, 3)
	g.Dexterity = NewCaracteristic("dexterity", 4, 6, 3)
	g.Constitution = NewCaracteristic("constitution", 4, 6, 3)

	g.CurrentHealth = 12 + g.Constitution.Modifier
	g.MaxHealth = 12 + g.Constitution.Modifier

	g.Weapon = NewRandomWeapon()
	g.Armor = NewRandomArmor()

	g.ArmorClass = calculateArmorClass(g)

	if level > 1 {
		for range dice.MakeRange(2, level) {
			g.LevelUp()
		}
		g.Experience = 0
	}

	return g
}

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

func (g *Gladiator) Attack() int {
	var result = dice.Roll(1, 20, -1)

	switch g.Weapon.Type {
	case MeleeWeapon:
		result += g.Strength.Modifier
	case RangedWeapon:
		result += g.Dexterity.Modifier
	}

	return result
}

func (g *Gladiator) Damage() int {
	var result = dice.Roll(g.Weapon.Number, g.Weapon.Damage, -1)

	switch g.Weapon.Type {
	case MeleeWeapon:
		result += g.Strength.Modifier
	case RangedWeapon:
		result += g.Dexterity.Modifier
	}

	return result
}
