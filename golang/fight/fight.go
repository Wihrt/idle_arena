package fight

import (
	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/gladiator"
	"go.uber.org/zap"
)

type FightResult struct {
	FightWon  bool                 `json:"fight_won"`
	Gladiator *gladiator.Gladiator `json:"gladiator"`
}

func ResolveFight(g *gladiator.Gladiator) (*FightResult, error) {
	var fightResult = &FightResult{
		FightWon:  false,
		Gladiator: g}

	fightWon, err := Fight(g)
	if err != nil {
		zap.L().Error("Error when generating fight",
			zap.Error(err),
		)
		return fightResult, err
	}

	fightResult.FightWon = fightWon
	g.CurrentHealth = g.MaxHealth
	if fightWon {
		expGained := dice.Roll(1, 20, -1)
		g.Experience += expGained
	}

	if g.Experience >= g.ExperienceToNextLevel {
		g.LevelUp()
	}

	return fightResult, nil
}

func Fight(player *gladiator.Gladiator) (bool, error) {
	enemy, err := gladiator.NewGladiator(player.Level, "")
	if err != nil {
		zap.L().Error("Error when creating enemy",
			zap.Error(err),
		)
		return false, err
	}

	zap.L().Info("Enemy created",
		zap.Int("strength", enemy.Strength.Value),
		zap.Int("dexterity", enemy.Dexterity.Value),
		zap.Int("constitution", enemy.Constitution.Value),
		zap.String("weapon", enemy.Weapon.Name),
		zap.String("armor", enemy.Armor.Name),
		zap.Int("armorClass", enemy.ArmorClass),
		zap.Int("currentHealth", enemy.CurrentHealth),
		zap.Int("maxHealth", enemy.MaxHealth),
		zap.Int("level", enemy.Level),
		zap.Int("experience", enemy.Experience),
		zap.Int("experienceNextLevel", enemy.ExperienceToNextLevel),
	)

	for {
		nextRound := Round(player, enemy)
		if !nextRound {
			break
		}
	}

	// Determine if the player has won the fight
	return player.CurrentHealth <= 0, nil
}

func Round(player *gladiator.Gladiator, enemy *gladiator.Gladiator) bool {
	var fightAgain = true

	if player.Attack() > enemy.ArmorClass {
		enemy.CurrentHealth -= player.Damage()
	}
	if enemy.Attack() > player.ArmorClass {
		player.CurrentHealth -= enemy.Damage()
	}

	if player.CurrentHealth <= 0 || enemy.CurrentHealth <= 0 {
		fightAgain = false
	}

	return fightAgain
}
