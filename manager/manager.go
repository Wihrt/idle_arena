package manager

import (
	"crypto/sha256"
	"encoding/hex"
)

type Manager struct {
	ManagerID  string   `json:"manager_id" bson:"manager_id"`
	Gladiators []string `json:"gladiators" bson:"gladiators"`
}

func NewManager(managerID string) *Manager {
	m := &Manager{
		ManagerID: managerID,
	}

	return m
}

func GenerateID(guildID string, userID string) string {
	h := sha256.Sum256([]byte(guildID + "." + userID))
	id := hex.EncodeToString(h[:])
	return id
}
