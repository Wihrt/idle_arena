package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/manager"
	"go.uber.org/zap"
)

func RegisterManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID  = generateManagerID(e)
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

func RetireManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
	var (
		mID  = generateManagerID(e)
		url  = os.Getenv("ARENA_URL")
		a    = arena.NewClient(url)
		data api.InteractionResponse
	)

	err := a.RetireManager(mID)
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

func generateManagerID(e *gateway.InteractionCreateEvent) string {
	id := manager.GenerateID(e.GuildID.String(), e.Member.User.ID.String())
	return id
}
