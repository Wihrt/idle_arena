package commands

import (
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/gladiator"
	"go.uber.org/zap"
)

func HireGladiator(e *gateway.InteractionCreateEvent) api.InteractionResponse {

	url := os.Getenv("ARENA_URL")
	a := arena.NewArenaClient(url)
	g, err := a.HireGladiator(e)
	if err != nil {
		zap.L().Error("Cannot hire a new gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
	}

	embed := GladiatorToEmbed(g)

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: []discord.Embed{embed},
		},
	}

	return data
}

func GetGladiator(e *gateway.InteractionCreateEvent) api.InteractionResponse {

	url := os.Getenv("ARENA_URL")
	a := arena.NewArenaClient(url)
	g, err := a.GetGladiator(e)
	if err != nil {
		zap.L().Error("Cannot get gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
	}

	embed := GladiatorToEmbed(g)

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: []discord.Embed{embed},
		},
	}

	return data
}

func FightGladiator(e *gateway.InteractionCreateEvent) api.InteractionResponse {

	url := os.Getenv("ARENA_URL")
	a := arena.NewArenaClient(url)
	g, err := a.FightGladiator(e)
	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
	}

	embed := GladiatorToEmbed(g)

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: []discord.Embed{embed},
		},
	}

	return data
}

func FireGladiator(e *gateway.InteractionCreateEvent) api.InteractionResponse {
	url := os.Getenv("ARENA_URL")
	a := arena.NewArenaClient(url)
	err := a.FireGladiator(e)
	if err != nil {
		zap.L().Error("Cannot fire gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
	}

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: "Gladiator fired",
		},
	}

	return data
}

func GladiatorToEmbed(g gladiator.Gladiator) discord.Embed {

	embed := discord.Embed{
		Title: g.Name,
		Fields: []discord.EmbedField{
			{Name: g.Strength.Name, Value: strconv.Itoa(g.Strength.Value)},
			{Name: g.Dexterity.Name, Value: strconv.Itoa(g.Dexterity.Value)},
			{Name: g.Constitution.Name, Value: strconv.Itoa(g.Constitution.Value)},
			{Name: "Weapon", Value: g.Weapon.Name},
			{Name: "Armor", Value: g.Armor.Name},
			{Name: "Armor Class", Value: strconv.Itoa(g.ArmorClass)},
			{Name: "Experience", Value: strconv.Itoa(g.Experience) + "/" + strconv.Itoa(g.ExperienceToNextLevel)},
		},
	}

	return embed
}
