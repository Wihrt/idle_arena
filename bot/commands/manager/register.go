package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func RegisterManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID  = utils.GenerateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		a    = arena.NewClient(url)
		data api.InteractionResponse
	)

	_, err := a.RegisterManager(mID)
	if err != nil {
		zap.L().Error("Cannot register manager",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
		return data, err
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString("You are now registered as an Arena Manager !"),
		},
	}

	return data, nil
}
