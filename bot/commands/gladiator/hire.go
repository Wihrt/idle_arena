package commands

import (
	"fmt"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"go.uber.org/zap"
)

func HireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "You have hired %s !"
		msg       string
		eArray    []discord.Embed
		mID       = utils.GenerateManagerID(e)
		url       = os.Getenv("ARENA_URL")
		a         = client.NewClient(url)
		data      api.InteractionResponse
	)

	g, err := a.HireGladiator(mID)
	if err != nil {
		zap.L().Error("Cannot hire a new gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
		return data, err
	}

	msg = fmt.Sprintf(formatMsg, g.Name)
	embed := GladiatorToEmbed(g)
	eArray = append(eArray, embed)

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(msg),
			Embeds:  &eArray,
		},
	}

	return data, nil
}
