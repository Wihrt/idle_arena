package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func FireGladiatorsMenu(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
	var (
		mID  = utils.GenerateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		c    = client.NewClient(url)
		data api.InteractionResponse
	)

	g, err := c.GetGladiators(mID)
	if err != nil {
		zap.L().Error("Cannot get gladiators",
			zap.String("managerID", mID),
			zap.Error(err),
		)
	}

	menu := GladiatorSelectMenu(g, "fire_gladiator_menu", 10)
	components := ComponentsWrapper([]discord.Component{menu})

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content:    option.NewNullableString("Select your gladiator to fire"),
			Components: &components,
		},
	}

	return data, nil
}

func FireGladiators(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID  = utils.GenerateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		c    = client.NewClient(url)
		data api.InteractionResponse
		msg  = "Your gladiators %s have been fired !"
	)

	for _, v := range e.Data.Values {
		gID := utils.GenerateGladiatorID(mID, v)
		err := c.FireGladiator(mID, gID)
		if err != nil {
			zap.L().Error("Cannot fire gladiator",
				zap.String("managerID", mID),
				zap.String("gladiatorID", gID),
				zap.String("gladiatorName", v),
				zap.Error(err),
			)
			return data, err
		}
	}

	msgFormatted := fmt.Sprintf(msg, strings.Join(e.Data.Values, ", "))

	data = api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(msgFormatted),
		},
	}

	return data, nil
}
