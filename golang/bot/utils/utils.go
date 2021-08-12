package utils

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/manager"
)

func GenerateManagerID(e *gateway.InteractionCreateEvent) string {
	id := manager.GenerateID(e.GuildID.String(), e.Member.User.ID.String())
	return id
}

func GenerateGladiatorID(managerID string, gladiatorName string) string {
	id := gladiator.GenerateID(managerID, gladiatorName)
	return id
}
