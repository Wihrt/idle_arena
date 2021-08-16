package commands

import (
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func RegisterManagerMenu(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
	easyButton := utils.Button("Easy", "easy_manager_difficulty", discord.SuccessButton)
	normalButton := utils.Button("Normal", "normal_manager_difficulty", discord.PrimaryButton)
	hardButton := utils.Button("Hard", "hard_manager_difficulty", discord.DangerButton)
	components := utils.ComponentsWrapper([]discord.Component{easyButton, normalButton, hardButton})

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content:    option.NewNullableString("Select the difficulty (easy, normal, hard)"),
			Components: &components,
		},
	}

	return data, nil
}

func RegisterManager(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID        = utils.GenerateManagerID(e)
		url        = os.Getenv("ARENA_URL")
		c          = client.NewClient(url)
		data       api.InteractionResponse
		difficulty int
	)

	switch e.Data.CustomID {
	case "easy_manager_difficulty":
		difficulty = 0
	case "normal_manager_difficulty":
		difficulty = 1
	case "hard_manager_difficulty":
		difficulty = 2
	}

	_, err := c.RegisterManager(mID, difficulty)
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
