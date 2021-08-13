package dnd

import (
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Weapon struct {
	Index          string       `json:"index" bson:"index"`
	Name           string       `json:"name" bson:"name"`
	WeaponCategory string       `json:"weapon_category" bson:"weapon_category"`
	WeaponRange    string       `json:"weapon_range" bson:"weapon_range"`
	CategoryRange  string       `json:"category_range" bson:"category_range"`
	Cost           Cost         `json:"cost" bson:"cost"`
	Damage         Damage       `json:"damage" bson:"damage"`
	Weight         float32      `json:"weight" bson:"weight"`
	Properties     []Properties `json:"properties" bson:"properties"`
	URL            string       `json:"url" bson:"url"`
}

type Damage struct {
	DamageDice string     `json:"damage_dice" bson:"damage_dice"`
	DamageType DamageType `json:"damage_type" bson:"damage_type"`
}

type DamageType struct {
	Index string `json:"index" bson:"index"`
	Name  string `json:"name" bson:"name"`
	URL   string `json:"url" bson:"url"`
}

type Range struct {
	Normal int `json:"normal" bson:"normal"`
	Long   int `json:"long" bson:"long"`
}

type Properties struct {
	Index string `json:"index" bson:"index"`
	Name  string `json:"name" bson:"name"`
	Url   string `json:"url" bson:"url"`
}

func (w *Weapon) HasFinesse() bool {
	for _, p := range w.Properties {
		if p.Index == "finesse" {
			return true
		}
	}
	return false
}

func (w *Weapon) ParseDice() ([]int, error) {

	var (
		parsedDice = strings.Split(w.Damage.DamageDice, "d")
		results    []int
	)

	for _, v := range parsedDice {
		i, err := strconv.Atoi(v)
		if err != nil {
			zap.L().Error("Cannot parse integer",
				zap.String("value", v),
				zap.Error(err),
			)
			return results, err
		}
		results = append(results, i)
	}
	return results, nil
}
