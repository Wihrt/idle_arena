package arena

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"go.uber.org/zap"
)

type Manager struct {
	GuildID    string   `json:"guild_id"`
	UserID     string   `json:"user_id"`
	Gladiators []string `json:"gladiators"`
}

func NewManager(guildID string, userID string) *Manager {
	m := &Manager{
		GuildID: guildID,
		UserID:  userID,
	}
	return m
}

func decodeManager(requestBody io.ReadCloser) *Manager {
	var manager Manager

	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		zap.L().Error("Cannot decode JSON")
	}
	json.Unmarshal(body, &manager)

	return &manager
}
