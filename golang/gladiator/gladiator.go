package gladiator

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/wihrt/idle_arena/dice"
	"go.uber.org/zap"
)

type Gladiator struct {
	ManagerID             string         `json:"manager_id" bson:"manager_id"`
	GladiatorID           string         `json:"gladiator_id" bson:"gladiator_id"`
	Armor                 *Armor         `json:"armor" bson:"armor"`
	ArmorClass            int            `json:"armor_class" bson:"armor_class"`
	Constitution          *Caracteristic `json:"constitution" bson:"constitution"`
	CurrentHealth         int            `json:"current_health" bson:"current_health"`
	Dexterity             *Caracteristic `json:"dexterity" bson:"dexterity"`
	Experience            int            `json:"experience" bson:"experience"`
	ExperienceToNextLevel int            `json:"experience_to_next_level" bson:"experience_to_next_level"`
	Level                 int            `json:"level" bson:"level"`
	MaxHealth             int            `json:"max_health" bson:"max_health"`
	Name                  string         `json:"name" bson:"name"`
	Strength              *Caracteristic `json:"strength" bson:"strength"`
	Weapon                *Weapon        `json:"weapon" bson:"weapon"`
}

func NewGladiator(level int, managerID string) (*Gladiator, error) {
	g := &Gladiator{
		Experience:            0,
		ExperienceToNextLevel: calculateNextLevel(1),
		Level:                 1,
		ManagerID:             managerID,
	}

	name, err := NewRandomName()
	if err != nil {
		zap.L().Error("Error when generating new name",
			zap.Error(err),
		)
		return g, err
	}
	g.Name = name

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

	g.generateID()

	return g, nil
}

func (g *Gladiator) generateID() {
	g.GladiatorID = GenerateID(g.ManagerID, g.Name)
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

func GenerateID(managerID string, name string) string {
	h := sha256.Sum256([]byte(managerID + "." + name))
	id := hex.EncodeToString(h[:])
	return id
}
