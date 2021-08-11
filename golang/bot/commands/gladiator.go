package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/gladiator"
	"go.uber.org/zap"
)

func HireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "You have hired %s !"
		msg       string
		mID       = generateManagerID(e)
		url       = os.Getenv("ARENA_URL")
		a         = arena.NewClient(url)
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
	embed := gladiatorToEmbed(g)

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: msg,
			Embeds:  []discord.Embed{embed},
		},
	}

	return data, nil
}

func GetGladiators(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		gArray []gladiator.Gladiator
		eArray []discord.Embed
		mID    = generateManagerID(e)
		name   = fetchValue(e.Data.Options, "name")
		gID    = generateGladiatorID(mID, name)
		url    = os.Getenv("ARENA_URL")
		a      = arena.NewClient(url)
		data   api.InteractionResponse
	)

	if len(e.Data.Options) == 1 {
		g, err := a.GetGladiator(mID, gID)
		if err != nil {
			zap.L().Error("Cannot get gladiators",
				zap.String("UserID", e.Member.User.ID.String()),
				zap.String("GuildID", e.GuildID.String()),
				zap.Error(err),
			)
			return data, err
		}
		gArray = append(gArray, g)
	} else {
		g, err := a.GetGladiators(mID)
		if err != nil {
			zap.L().Error("Cannot get gladiator",
				zap.String("UserID", e.Member.User.ID.String()),
				zap.String("GuildID", e.GuildID.String()),
				zap.Error(err),
			)
			return data, err
		}
		gArray = append(gArray, g...)
	}

	for _, g := range gArray {
		e := gladiatorToEmbed(g)
		eArray = append(eArray, e)
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: eArray,
		},
	}

	return data, nil
}

func FightGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "Your gladiator %s has %s the fight !"
		msg       string
		mID       = generateManagerID(e)
		name      = fetchValue(e.Data.Options, "name")
		gID       = generateGladiatorID(mID, name)
		url       = os.Getenv("ARENA_URL")
		a         = arena.NewClient(url)
		data      api.InteractionResponse
	)

	f, err := a.FightGladiator(mID, gID)
	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
		return data, err
	}

	if f.FightWon {
		msg = fmt.Sprintf(formatMsg, f.Gladiator.Name, "won")
	} else {
		msg = fmt.Sprintf(formatMsg, f.Gladiator.Name, "lost")
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: msg,
		},
	}

	return data, nil
}

func FireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "You have fired %s !"
		msg       string
		mID       = generateManagerID(e)
		name      = fetchValue(e.Data.Options, "name")
		gID       = generateGladiatorID(mID, name)
		url       = os.Getenv("ARENA_URL")
		a         = arena.NewClient(url)
		data      api.InteractionResponse
	)
	err := a.FireGladiator(mID, gID)
	if err != nil {
		zap.L().Error("Cannot fire gladiator",
			zap.String("UserID", e.Member.User.ID.String()),
			zap.String("GuildID", e.GuildID.String()),
			zap.Error(err),
		)
		return data, err
	}

	msg = fmt.Sprintf(formatMsg, name)

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: msg,
		},
	}

	return data, err
}

func gladiatorToEmbed(g gladiator.Gladiator) discord.Embed {

	embed := discord.Embed{
		Title: g.Name,
		Fields: []discord.EmbedField{
			{Name: g.Strength.Name, Value: strconv.Itoa(g.Strength.Value), Inline: true},
			{Name: g.Dexterity.Name, Value: strconv.Itoa(g.Dexterity.Value), Inline: true},
			{Name: g.Constitution.Name, Value: strconv.Itoa(g.Constitution.Value), Inline: true},
			{Name: "Weapon", Value: g.Weapon.Name, Inline: true},
			{Name: "Armor", Value: g.Armor.Name, Inline: true},
			{Name: "Armor Class", Value: strconv.Itoa(g.ArmorClass), Inline: true},
			{Name: "Experience", Value: strconv.Itoa(g.Experience) + "/" + strconv.Itoa(g.ExperienceToNextLevel), Inline: false},
		},
	}

	return embed
}

func generateGladiatorID(managerID string, gladiatorName string) string {
	id := gladiator.GenerateID(managerID, gladiatorName)
	return id
}
