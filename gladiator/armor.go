package gladiator

import (
	"context"
	"time"

	"github.com/wihrt/idle_arena/dnd"
	"github.com/wihrt/idle_arena/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func NewRandomArmor(m *mongo.Client) (*dnd.Armor, error) {

	var (
		a           dnd.Armor
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	sampleStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	matchStage := bson.D{{"$match", bson.D{{"str_minimum", 0}}}}

	cur, err := m.Database(utils.DB).Collection(utils.A).Aggregate(ctx, mongo.Pipeline{matchStage, sampleStage})
	if err != nil {
		zap.L().Error("Cannot get armor",
			zap.String("database", utils.DB),
			zap.String("collection", utils.A),
			zap.Error(err),
		)
		return &a, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err = cur.Decode(&a)
		if err != nil {
			zap.L().Error("Cannot decode armor",
				zap.Error(err),
			)
			return &a, err
		}
	}

	return &a, nil
}

func calculateArmorClass(g *Gladiator) int {
	var (
		armorClass = g.Armor.ArmorClass.Base
		dexBonus   = g.Dexterity.Modifier
	)

	if g.Armor.ArmorClass.DexBonus {
		if g.Armor.ArmorClass.MaxBonus != 0 {
			if dexBonus > g.Armor.ArmorClass.MaxBonus {
				armorClass += g.Armor.ArmorClass.MaxBonus
			} else {
				armorClass += dexBonus
			}
		}
	}

	return armorClass
}
