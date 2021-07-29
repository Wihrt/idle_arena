package commands

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"go.uber.org/zap"
)

var toto = discord.CommandOption{
	Type:        0,
	Name:        "name",
	Description: "",
	Required:    false,
	Choices:     []discord.CommandOptionChoice{},
	Options:     []discord.CommandOption{},
}

var RegisteredCommands = []api.CreateCommandData{
	{
		Name:        "register",
		Description: "Register yourself as a new Arena Manager",
		Options:     []discord.CommandOption{},
	},
	{
		Name:        "retire",
		Description: "Retire yourself as an Arena Manager",
	},
	{
		Name:        "hire",
		Description: "Hire a new Gladiator",
	},
	{
		Name:        "show",
		Description: "Show your gladiator",
		Options: []discord.CommandOption{
			{
				Type:        3,
				Name:        "name",
				Description: "Name of your gladiator",
				Required:    false,
			},
		},
	},
	{
		Name:        "fight",
		Description: "Make your gladiator performs a fight",
		Options: []discord.CommandOption{
			{
				Type:        3,
				Name:        "name",
				Description: "Name of your gladiator",
				Required:    true,
			},
		},
	},
	{
		Name:        "fire",
		Description: "Fire your gladiator",
		Options: []discord.CommandOption{
			{
				Type:        3,
				Name:        "name",
				Description: "Name of your gladiator",
				Required:    true,
			},
		},
	},
}

func HandleInteraction(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		data api.InteractionResponse
		err  error
	)

	switch e.Data.Name {
	case "register":
		data, err = RegisterManager(e)
	case "retire":
		data, err = RetireManager(e)
	case "hire":
		data, err = HireGladiator(e)
	case "show":
		data, err = GetGladiators(e)
	case "fight":
		data, err = FightGladiator(e)
	case "fire":
		data, err = FireGladiator(e)
	}

	if err != nil {
		zap.L().Error("Error when invoking command",
			zap.String("commandName", e.Data.Name),
			zap.Error(err),
		)
		return data, err
	}

	return data, nil
}

func fetchValue(options []gateway.InteractionOption, optionName string) string {
	var value string

	for _, o := range options {
		if o.Name == optionName {
			value = o.Value
		}
	}

	zap.L().Debug("Result of fetch value",
		zap.String("name", optionName),
		zap.String("value", value),
	)

	return value
}
