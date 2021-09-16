package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func ShowManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID    = utils.GenerateManagerID(e)
		url    = os.Getenv("ARENA_URL")
		c      = client.NewClient(url)
		data   api.InteractionResponse
		eArray []discord.Embed
	)

	m, err := c.ShowManager(mID)
	if err != nil {
		zap.L().Error("Cannot get manager",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return data, err
	}

	g, err := c.GetGladiators(mID)
	if err != nil {
		zap.L().Error("Cannot get gladiators",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return data, err
	}

	embed := utils.ManagerToEmbed(*m, g)
	eArray = append(eArray, embed)

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &eArray,
		},
	}

	return data, nil
}
