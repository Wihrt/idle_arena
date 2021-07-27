package commands

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

var RegisteredCommands = []api.CreateCommandData{
	{
		Name:        "ping",
		Description: "Basic ping command.",
		Options:     []discord.CommandOption{},
	},
	{
		Name:        "hire",
		Description: "Hire a new Gladiator",
	},
	{
		Name:        "show",
		Description: "Show your gladiator",
	},
	{
		Name:        "fight",
		Description: "Make your gladiator performs a fight",
	},
	{
		Name:        "fire",
		Description: "Fire your gladiator",
	},
}

func HandleInteraction(e *gateway.InteractionCreateEvent) api.InteractionResponse {

	var data api.InteractionResponse

	switch e.Data.Name {
	case "ping":
		data = PingCommand(e)
	case "hire":
		data = HireGladiator(e)
	case "show":
		data = GetGladiator(e)
	case "fight":
		data = FightGladiator(e)
	case "fire":
		data = FireGladiator(e)
	}

	return data
}
