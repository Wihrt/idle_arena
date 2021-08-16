package commands

import (
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/arena/errors"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func HireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		embed  discord.Embed
		eArray []discord.Embed
		mID    = utils.GenerateManagerID(e)
		url    = os.Getenv("ARENA_URL")
		a      = client.NewClient(url)
		data   api.InteractionResponse
		number = 1
	)

	for _, o := range e.Data.Options {
		switch o.Name {
		case "number":
			number, _ = strconv.Atoi(o.Value.String())
		}
	}

	seq := make([]int, number)
	for range seq {
		g, err := a.HireGladiator(mID)
		if err != nil && err != errors.ErrNotEnoughMoney {
			zap.L().Error("Cannot hire a new gladiator",
				zap.String("UserID", e.Member.User.ID.String()),
				zap.String("GuildID", e.GuildID.String()),
				zap.Error(err),
			)
			return data, err
		}

		if err != nil && err == errors.ErrNotEnoughMoney {
			embed = discord.Embed{
				Title: "You don't have enough money !",
				Type:  discord.NormalEmbed,
				Color: 0xff0000,
			}

		} else {
			embed = utils.GladiatorToEmbed(g)
			embed.Color = 0x00ff00
		}
		eArray = append(eArray, embed)
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &eArray,
		},
	}

	return data, nil
}
