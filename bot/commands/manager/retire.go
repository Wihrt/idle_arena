package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func RetireManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
	var (
		mID  = utils.GenerateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		c    = client.NewClient(url)
		data api.InteractionResponse
	)

	err := c.RetireManager(mID)
	if err != nil {
		zap.L().Error("Cannot retire manager",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
		return data, err
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString("You are retired as an Arena Manager !"),
		},
	}

	return data, nil
}
