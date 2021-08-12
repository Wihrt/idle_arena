package gladiator

import (
	"context"
	"time"

	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/dnd"
	"github.com/wihrt/idle_arena/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewRandomWeapon(m *mongo.Client) (*dnd.Weapon, error) {
	var (
		w           dnd.Weapon
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	sampleStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cur, err := m.Database(utils.DB).Collection(utils.W).Aggregate(ctx, mongo.Pipeline{sampleStage})
	if err != nil {
		zap.L().Error("Cannot get weapon",
			zap.String("database", utils.DB),
			zap.String("collection", utils.W),
			zap.Error(err),
		)
		return &w, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err = cur.Decode(&w)
		if err != nil {
			zap.L().Error("Cannot decode weapon",
				zap.Error(err),
			)
			return &w, err
		}
	}

	return &w, nil
}

func (g *Gladiator) Attack() int {
	var (
		result   = dice.Roll(1, 20, -1)
		modifier int
	)

	switch g.Weapon.WeaponRange {
	case "Melee":
		modifier = g.Strength.Modifier
		if g.Weapon.HasFinesse() {
			modifier = g.Dexterity.Modifier
		}
	case "Ranged":
		modifier = g.Dexterity.Modifier
	}

	result += modifier
	return result
}

func (g *Gladiator) Damage() int {

	var modifier int

	weaponDice, _ := g.Weapon.ParseDice()
	result := dice.Roll(weaponDice[0], weaponDice[1], -1)

	switch g.Weapon.WeaponRange {
	case "Melee":
		modifier = g.Strength.Modifier
		if g.Weapon.HasFinesse() {
			modifier = g.Dexterity.Modifier
		}
	case "Ranged":
		modifier = g.Dexterity.Modifier
	}

	result += modifier
	return result
}
