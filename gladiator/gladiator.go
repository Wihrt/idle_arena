package gladiator

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/dnd"
	"github.com/wihrt/idle_arena/manager"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Gladiator struct {
	ManagerID    string         `json:"manager_id" bson:"manager_id"`
	GladiatorID  string         `json:"gladiator_id" bson:"gladiator_id"`
	Armor        *dnd.Armor     `json:"armor" bson:"armor"`
	ArmorClass   int            `json:"armor_class" bson:"armor_class"`
	Constitution *Caracteristic `json:"constitution" bson:"constitution"`
	Dexterity    *Caracteristic `json:"dexterity" bson:"dexterity"`
	Experience   *Experience    `json:"experience" bson:"experience"`
	Health       *Health        `json:"health" bson:"health"`
	Name         string         `json:"name" bson:"name"`
	Strength     *Caracteristic `json:"strength" bson:"strength"`
	Weapon       *dnd.Weapon    `json:"weapon" bson:"weapon"`
	DeathSave    *DeathSave     `json:"death_saves" bson:"death_saves"`
}

func NewGladiator(level int, m *manager.Manager, mongoClient *mongo.Client) (*Gladiator, error) {
	g := &Gladiator{
		ManagerID:  m.ManagerID,
		Experience: NewExperience(),
		DeathSave:  NewDeathSave(m.Difficulty),
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
	g.Health = NewHealth(g.Constitution.Modifier)

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
		g.Experience.Current = 0
	}

	g.GladiatorID = GenerateID(m.ManagerID, name)

	return g, nil
}

func NewEnemy(level int, mongoClient *mongo.Client) (*Gladiator, error) {
	g := &Gladiator{
		Experience: NewExperience(),
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
	g.Health = NewHealth(g.Constitution.Modifier)

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
		g.Experience.Current = 0
	}

	return g, nil
}

func GenerateID(managerID string, name string) string {
	h := sha256.Sum256([]byte(managerID + "." + name))
	id := hex.EncodeToString(h[:])
	return id
}
