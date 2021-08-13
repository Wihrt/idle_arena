package gladiator

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/dnd"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Gladiator struct {
	ManagerID             string         `json:"manager_id" bson:"manager_id"`
	GladiatorID           string         `json:"gladiator_id" bson:"gladiator_id"`
	Armor                 *dnd.Armor     `json:"armor" bson:"armor"`
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
	Weapon                *dnd.Weapon    `json:"weapon" bson:"weapon"`
}

func NewGladiator(level int, managerID string, mongoClient *mongo.Client) (*Gladiator, error) {
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

	g.Weapon, err = NewRandomWeapon(mongoClient)
	if err != nil {
		zap.L().Error("Error when generating new weapon",
			zap.Error(err),
		)
		return g, err
	}
	g.Armor, err = NewRandomArmor(mongoClient)
	if err != nil {
		zap.L().Error("Error when generating new armor",
			zap.Error(err),
		)
		return g, err
	}

	g.ArmorClass = calculateArmorClass(g)

	if level > 1 {
		for range dice.MakeRange(2, level) {
			g.LevelUp()
		}
		g.Experience = 0
	}

	g.GladiatorID = GenerateID(managerID, name)

	return g, nil
}

func GenerateID(managerID string, name string) string {
	h := sha256.Sum256([]byte(managerID + "." + name))
	id := hex.EncodeToString(h[:])
	return id
}
