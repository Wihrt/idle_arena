package commands

import (
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	gcommands "github.com/wihrt/idle_arena/bot/commands/gladiator"
	mcommands "github.com/wihrt/idle_arena/bot/commands/manager"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var RegisteredCommands = []api.CreateCommandData{
	{
		Name:        "register",
		Description: "Register yourself as a new Arena Manager",
	},
	{
		Name:        "retire",
		Description: "Retire yourself as an Arena Manager",
	},
	{
		Name:        "hire",
		Description: "Hire a new Gladiator",
		Options: []discord.CommandOption{
			{
				Type:        discord.IntegerOption,
				Name:        "number",
				Description: "Number of gladiator you want to hire",
				Required:    false,
			},
		},
	},
	{
		Name:        "show",
		Description: "Show your manager/gladiators",
		Options: []discord.CommandOption{
			{
				Type:        discord.StringOption,
				Name:        "type",
				Description: "Select what you want to show",
				Required:    true,
				Choices: []discord.CommandOptionChoice{
					{
						Name:  "Manager",
						Value: "manager",
					},
					{
						Name:  "Gladiators",
						Value: "gladiators",
					},
				},
			},
		},
	},
	{
		Name:        "fight",
		Description: "Make your gladiator performs a fight",
		Options: []discord.CommandOption{
			{
				Type:        discord.StringOption,
				Name:        "difficulty",
				Description: "Difficulty of the fight",
				Required:    true,
				Choices: []discord.CommandOptionChoice{
					{
						Name:  "Easy",
						Value: "easy",
					},
					{
						Name:  "Normal",
						Value: "normal",
					},
					{
						Name:  "Hard",
						Value: "hard",
					},
					{
						Name:  "Challenging",
						Value: "challenging",
					},
					{
						Name:  "Nightmarish",
						Value: "nightmarish",
					},
				},
			},
		},
	},
	{
		Name:        "heal",
		Description: "Heal your gladiator (reset death saves to 0)",
	},
	{
		Name:        "fire",
		Description: "Fire your gladiator(s)",
	},
}

func HandleInteraction(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		data api.InteractionResponse
		err  error
	)

	var fields = []zapcore.Field{
		zap.String("name", e.Data.Name),
		zap.String("custom_id", e.Data.CustomID),
	}
	for i, o := range e.Data.Options {
		fields = append(fields,
			zap.String("Option "+strconv.Itoa(i+1)+" name", o.Name),
			zap.String("Option "+strconv.Itoa(i+1)+" value", o.Value.String()))
	}
	zap.L().Debug("Command name",
		fields...,
	)

	// Here are the commands sent the first time by the user
	switch e.Data.Name {
	case "register":
		data, err = mcommands.RegisterManagerMenu(e)
	case "retire":
		data, err = mcommands.RetireManager(e)
	case "hire":
		data, err = gcommands.HireGladiator(e)
	case "show":
		value, _ := strconv.Unquote(e.Data.Options[0].Value.String())
		switch value {
		case "manager":
			data, err = mcommands.ShowManager(e)
		case "gladiators":
			data, err = gcommands.ShowGladiatorsMenu(e)
		}
	case "fight":
		data, err = gcommands.FightGladiatorsMenu(e)
	case "fire":
		data, err = gcommands.FireGladiatorsMenu(e)
	case "heal":
		data, err = gcommands.HealGladiatorsMenu(e)
	}

	// Here are the Component Interaction Response
	switch e.Data.CustomID {
	case "easy_manager_difficulty",
		"normal_manager_difficulty",
		"hard_manager_difficulty":
		data, err = mcommands.RegisterManager(e)
	case "show_gladiator_menu":
		data, err = gcommands.ShowGladiators(e)
	case "easy_fight_gladiator_menu",
		"normal_fight_gladiator_menu",
		"hard_fight_gladiator_menu",
		"challenging_fight_gladiator_menu",
		"nightmarish_fight_gladiator_menu",
		"hellish_fight_gladiator_menu":
		data, err = gcommands.FightGladiator(e)
	case "fire_gladiator_menu":
		data, err = gcommands.FireGladiators(e)
	case "heal_gladiator_menu":
		data, err = gcommands.HealGladiators(e)
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
