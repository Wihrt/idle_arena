package gladiator

import "github.com/wihrt/idle_arena/dice"

type ArmorType string

const (
	NoArmor     ArmorType = "none"
	LightArmor  ArmorType = "light"
	MediumArmor ArmorType = "medium"
	HeavyArmor  ArmorType = "heavy"
)

type Armor struct {
	Name        string    `json:"name"`
	Value       int       `json:"value"`
	Type        ArmorType `json:"type"`
	MaxDexBonus int       `json:"max_dex_bonus"`
}

func NewArmor(name string, value int, armorType ArmorType) *Armor {
	var a = &Armor{
		Name:        name,
		Value:       value,
		Type:        armorType,
		MaxDexBonus: 0,
	}

	a.calculateMaxBonus()
	return a
}

func NewRandomArmor() *Armor {
	var armor *Armor
	roll := dice.Roll(1, 4, -1)

	switch roll {
	case 1:
		armor = NewArmor("Slip", 0, NoArmor)
	case 2:
		armor = NewArmor("Leather Armor", 11, LightArmor)
	case 3:
		armor = NewArmor("Chain Shirt", 14, MediumArmor)
	case 4:
		armor = NewArmor("Chain Mail", 16, HeavyArmor)
	}

	return armor
}

func (a *Armor) calculateMaxBonus() {
	switch a.Type {
	case LightArmor:
		a.MaxDexBonus = -1
	case MediumArmor:
		a.MaxDexBonus = 2
	case HeavyArmor:
		a.MaxDexBonus = 0
	}
}

func calculateArmorClass(g *Gladiator) int {
	var armorClass int

	switch g.Armor.Type {
	case NoArmor:
		armorClass = 10 + g.Dexterity.Modifier
	case LightArmor:
		armorClass = g.Armor.Value + g.Dexterity.Modifier
	case MediumArmor:
		dexBonus := g.Dexterity.Modifier
		if dexBonus > g.Armor.MaxDexBonus {
			dexBonus = g.Armor.MaxDexBonus
		}
		armorClass = g.Armor.Value + dexBonus
	case HeavyArmor:
		armorClass = g.Armor.Value
	}

	return armorClass
}
