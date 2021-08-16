package fight

import (
	"github.com/wihrt/idle_arena/dice"
	"github.com/wihrt/idle_arena/gladiator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Result struct {
	FightWon       bool                 `json:"fight_won"`
	KilledInCombat bool                 `json:"killed_in_combat"`
	Gladiator      *gladiator.Gladiator `json:"gladiator"`
	Enemy          *gladiator.Gladiator `json:"enemy"`
}

func ResolveFight(g *gladiator.Gladiator, m *mongo.Client, s *Settings) (*Result, error) {
	var fightResult = &Result{
		FightWon:       false,
		KilledInCombat: false,
		Gladiator:      g}

	fightWon, enemy, err := Fight(g, m, s)
	if err != nil {
		zap.L().Error("Error when generating fight",
			zap.Error(err),
		)
		return fightResult, err
	}

	if !fightWon {
		deathSave := dice.Roll(1, 20, -1)
		zap.L().Debug("Result of death save",
			zap.Int("result", deathSave),
			zap.Bool("failed", deathSave < 10),
		)
		if deathSave < 10 {
			g.CurrentDeathSaves += 1
		}
	}

	if g.CurrentDeathSaves == g.MaxDeathSaves {
		fightResult.KilledInCombat = true
	}

	fightResult.FightWon = fightWon
	fightResult.Enemy = enemy
	g.CurrentHealth = g.MaxHealth
	if fightWon {
		expGained := dice.Roll(int(s.Difficulty)+1, 20, -1)
		g.Experience += expGained
	}

	if g.Experience >= g.ExperienceToNextLevel {
		g.LevelUp()
	}

	return fightResult, nil
}

func Fight(player *gladiator.Gladiator, m *mongo.Client, s *Settings) (bool, *gladiator.Gladiator, error) {
	var (
		enemy *gladiator.Gladiator
	)

	enemy, err := gladiator.NewEnemy(player.Level+int(s.Difficulty), m)
	if err != nil {
		zap.L().Error("Error when creating enemy",
			zap.Error(err),
		)
		return false, enemy, err
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
	return !(player.CurrentHealth <= 0), enemy, nil
}

func Round(player *gladiator.Gladiator, enemy *gladiator.Gladiator) bool {
	var fightAgain = true

	pAttack := player.Attack()
	zap.L().Debug("Player attacks enemy",
		zap.Int("Attack roll", pAttack),
		zap.Int("Armor class", enemy.ArmorClass),
		zap.Bool("Hit", pAttack > enemy.ArmorClass),
	)

	if pAttack > enemy.ArmorClass {
		pDamage := player.Damage()
		zap.L().Debug("Player damages enemy",
			zap.Int("Damage roll", pDamage),
			zap.Int("Enemy health", enemy.CurrentHealth),
			zap.Int("Enemy after hit", enemy.CurrentHealth-pDamage),
		)
		enemy.CurrentHealth -= pDamage
	}

	eAttack := enemy.Attack()
	zap.L().Debug("Enemy attacks player",
		zap.Int("Attack roll", eAttack),
		zap.Int("Armor class", player.ArmorClass),
		zap.Bool("Hit", eAttack > player.ArmorClass),
	)

	if eAttack > player.ArmorClass {
		eDamage := enemy.Damage()
		zap.L().Debug("Enemy damages player",
			zap.Int("Damage roll", eDamage),
			zap.Int("Player health", player.CurrentHealth),
			zap.Int("Player after hit", player.CurrentHealth-eDamage),
		)
		player.CurrentHealth -= eDamage
	}

	zap.L().Debug("End of round",
		zap.Int("Player health", player.CurrentHealth),
		zap.Int("Enemy health", enemy.CurrentHealth),
	)

	if player.CurrentHealth <= 0 || enemy.CurrentHealth <= 0 {
		fightAgain = false
	}

	return fightAgain
}
