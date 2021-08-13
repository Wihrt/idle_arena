package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func ShowGladiatorsMenu(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
	var (
		mID  = utils.GenerateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		a    = arena.NewClient(url)
		data api.InteractionResponse
	)

	g, err := a.GetGladiators(mID)
	if err != nil {
		zap.L().Error("Cannot get gladiators",
			zap.String("managerID", mID),
			zap.Error(err),
		)
	}

	menu := GladiatorSelectMenu(g, "show_gladiator_menu", 10)
	components := ComponentsWrapper([]discord.Component{menu})

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content:    option.NewNullableString("Select your gladiator to show"),
			Components: &components,
		},
	}

	return data, nil
}

func ShowGladiators(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID    = utils.GenerateManagerID(e)
		url    = os.Getenv("ARENA_URL")
		a      = arena.NewClient(url)
		data   api.InteractionResponse
		eArray []discord.Embed
	)

	for _, v := range e.Data.Values {
		gID := utils.GenerateGladiatorID(mID, v)
		g, err := a.GetGladiator(mID, gID)
		if err != nil {
			zap.L().Error("Cannot get gladiators",
				zap.String("managerID", mID),
				zap.String("gladiatorID", gID),
				zap.Error(err),
			)
			return data, err
		}

		e := GladiatorToEmbed(g)
		eArray = append(eArray, e)

	}

	data = api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Embeds: &eArray,
		},
	}

	return data, nil
}
