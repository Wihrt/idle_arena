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
	},
	{
		Name:        "show",
		Description: "Show your gladiator(s)",
	},
	{
		Name:        "fight",
		Description: "Make your gladiator performs a fight",
		Options: []discord.CommandOption{
			{
				Type:        discord.SubcommandOption,
				Name:        "easy",
				Description: "Fight an easy enemy",
			},
			{
				Type:        discord.SubcommandOption,
				Name:        "normal",
				Description: "Fight an normal enemy",
			},
			{
				Type:        discord.SubcommandOption,
				Name:        "hard",
				Description: "Fight a difficult enemy",
			},
			{
				Type:        discord.SubcommandOption,
				Name:        "challenging",
				Description: "Fight a challenging enemy",
			},
			{
				Type:        discord.SubcommandOption,
				Name:        "nightmarish",
				Description: "Fight a nightmarish enemy",
			},
			{
				Type:        discord.SubcommandOption,
				Name:        "hellish",
				Description: "Fight a hellish enemy",
			},
		},
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
		data, err = mcommands.RegisterManager(e)
	case "retire":
		data, err = mcommands.RetireManager(e)
	case "hire":
		data, err = gcommands.HireGladiator(e)
	case "show":
		data, err = gcommands.ShowGladiatorsMenu(e)
	case "fight":
		data, err = gcommands.FightGladiatorsMenu(e)
	case "fire":
		data, err = gcommands.FireGladiatorsMenu(e)
	}

	// Here are the Component Interaction Response
	switch e.Data.CustomID {
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
