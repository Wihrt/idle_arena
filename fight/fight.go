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
			g.DeathSave.Current = g.DeathSave.Current + 1
		}
	}

	if g.DeathSave.Current == g.DeathSave.Max {
		fightResult.KilledInCombat = true
	}

	fightResult.FightWon = fightWon
	fightResult.Enemy = enemy
	g.Health.Current = g.Health.Max
	if fightWon {
		expGained := dice.Roll(int(s.Difficulty)+1, 20, -1)
		g.Experience.Current += expGained
	}

	if g.Experience.Current >= g.Experience.NextLevel {
		g.LevelUp()
	}

	return fightResult, nil
}

func Fight(player *gladiator.Gladiator, m *mongo.Client, s *Settings) (bool, *gladiator.Gladiator, error) {
	var (
		enemy *gladiator.Gladiator
	)

	enemy, err := gladiator.NewEnemy(player.Experience.Level+int(s.Difficulty), m)
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
		zap.Int("currentHealth", enemy.Health.Current),
		zap.Int("maxHealth", enemy.Health.Max),
		zap.Int("level", enemy.Experience.Level),
		zap.Int("experience", enemy.Experience.Current),
		zap.Int("experienceNextLevel", enemy.Experience.NextLevel),
	)

	for {
		nextRound := Round(player, enemy)
		if !nextRound {
			break
		}
	}

	// Determine if the player has won the fight
	return !(player.Health.Current <= 0), enemy, nil
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
			zap.Int("Enemy health", enemy.Health.Current),
			zap.Int("Enemy after hit", enemy.Health.Current-pDamage),
		)
		enemy.Health.Current -= pDamage
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
			zap.Int("Player health", player.Health.Current),
			zap.Int("Player after hit", player.Health.Current-eDamage),
		)
		player.Health.Current -= eDamage
	}

	zap.L().Debug("End of round",
		zap.Int("Player health", player.Health.Current),
		zap.Int("Enemy health", enemy.Health.Current),
	)

	if player.Health.Current <= 0 || enemy.Health.Current <= 0 {
		fightAgain = false
	}

	return fightAgain
}
