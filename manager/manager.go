package manager

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Manager struct {
	ManagerID  string          `json:"manager_id" bson:"manager_id"`
	Name       string          `json:"name" bson:"name"`
	GuildID    discord.GuildID `json:"guild_id" bson:"guild_id"`
	Difficulty Difficulty      `json:"difficulty" bson:"difficulty"`
	MoneyPouch *MoneyPouch     `json:"money_pouch"`
	Gladiators []string        `json:"gladiators" bson:"gladiators"`
}

func NewManager(managerID string, name string, guildID discord.GuildID, difficulty int) (*Manager, error) {
	m := &Manager{
		Name:      name,
		GuildID:   guildID,
		ManagerID: managerID,
	}

	d, err := ParseDifficulty(difficulty)
	if err != nil {
		return m, err
	}

	m.Difficulty = d
	m.MoneyPouch = NewMoneyPouch(d)

	return m, nil
}

func GenerateID(guildID string, userID string) string {
	h := sha256.Sum256([]byte(guildID + "." + userID))
	id := hex.EncodeToString(h[:])
	return id
}
