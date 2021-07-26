package fight

import (
	"github.com/wihrt/idle_arena/arena/dice"
	"github.com/wihrt/idle_arena/arena/gladiator"
	"go.uber.org/zap"
)

func ResolveFight(player *gladiator.Gladiator) {
	fightWon := Fight(player)
	player.CurrentHealth = player.MaxHealth
	if fightWon {
		expGained := dice.Roll(1, 20, -1)
		player.Experience += expGained
	}

	if player.Experience >= player.ExperienceToNextLevel {
		player.LevelUp()
	}
}

func Fight(player *gladiator.Gladiator) bool {
	var enemy = gladiator.NewGladiator(player.Level)

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
	return player.CurrentHealth <= 0
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
