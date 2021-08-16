package commands

import (
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena/client"
	"github.com/wihrt/idle_arena/bot/utils"
	"github.com/wihrt/idle_arena/fight"
	"go.uber.org/zap"
)

func FightGladiatorsMenu(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {
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

	menu := utils.GladiatorSelectMenu(g, e.Data.Options[0].Name+"_fight_gladiator_menu", 1)
	components := utils.ComponentsWrapper([]discord.Component{menu})

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content:    option.NewNullableString("Select your gladiator to fight"),
			Components: &components,
		},
	}

	return data, nil
}

func FightGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		mID        = utils.GenerateManagerID(e)
		url        = os.Getenv("ARENA_URL")
		a          = client.NewClient(url)
		data       api.InteractionResponse
		eArray     []discord.Embed
		difficulty = strings.Split(e.Data.CustomID, "_")[0]
	)

	s, err := fight.NewSettings(difficulty)
	if err != nil {
		zap.L().Error("Cannot generate fight settings",
			zap.String("difficulty", difficulty),
			zap.Error(err),
		)
		return data, err
	}

	for _, v := range e.Data.Values {
		gID := utils.GenerateGladiatorID(mID, v)
		f, err := a.FightGladiator(mID, gID, s)
		if err != nil {
			zap.L().Error("Cannot fight gladiator",
				zap.String("UserID", e.Member.User.ID.String()),
				zap.String("GuildID", e.GuildID.String()),
				zap.Error(err),
			)
			return data, err
		}
		embed := utils.FightToEmbed(f)
		eArray = append(eArray, embed)

		if f.KilledInCombat {
			err := a.FireGladiator(mID, gID)
			if err != nil {
				zap.L().Error("Cannot fire gladiator",
					zap.String("UserID", e.Member.User.ID.String()),
					zap.String("GuildID", e.GuildID.String()),
					zap.String("GladiatorID", gID),
					zap.String("Manager ID", mID),
					zap.Error(err),
				)
			}
		}
	}

	data = api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Embeds: &eArray,
		},
	}

	return data, nil
}
