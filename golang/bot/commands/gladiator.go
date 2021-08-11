package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/utils"
	"go.uber.org/zap"
)

func HireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "You have hired %s !"
		msg       string
		eArray    []discord.Embed
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

func GetGladiators(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		gArray []gladiator.Gladiator
		eArray []discord.Embed
		mID    = generateManagerID(e)
		// name   = fetchValue(e.Data.Options, "name")
		// gID    = generateGladiatorID(mID, name)
		url         = os.Getenv("ARENA_URL")
		a           = arena.NewClient(url)
		data        api.InteractionResponse
		menuOptions []discord.SelectComponentOption
	)

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

	for _, g := range gArray {
		option := discord.SelectComponentOption{
			Label: g.Name,
			Value: g.Name,
		}
		menuOptions = append(menuOptions, option)
	}
	allOptions := discord.SelectComponentOption{
		Label:   "all",
		Value:   "all",
		Default: true,
	}
	menuOptions = append(menuOptions, allOptions)

	selectMenu := discord.SelectComponent{
		CustomID: "name",
		Options:  menuOptions,
		Disabled: false,
	}

	// if len(e.Data.Options) == 1 {
	// 	g, err := a.GetGladiator(mID, gID)
	// 	if err != nil {
	// 		zap.L().Error("Cannot get gladiators",
	// 			zap.String("UserID", e.Member.User.ID.String()),
	// 			zap.String("GuildID", e.GuildID.String()),
	// 			zap.Error(err),
	// 		)
	// 		return data, err
	// 	}
	// 	gArray = append(gArray, g)
	// } else {
	// 	g, err := a.GetGladiators(mID)
	// 	if err != nil {
	// 		zap.L().Error("Cannot get gladiator",
	// 			zap.String("UserID", e.Member.User.ID.String()),
	// 			zap.String("GuildID", e.GuildID.String()),
	// 			zap.Error(err),
	// 		)
	// 		return data, err
	// 	}
	// 	gArray = append(gArray, g...)
	// }

	for _, g := range gArray {
		e := gladiatorToEmbed(g)
		eArray = append(eArray, e)
	}

	data = api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Components: &[]discord.Component{selectMenu},
		},
	}

	return data, nil
}

func FightGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "Your gladiator %s has %s the fight !"
		msg       string
		mID       = generateManagerID(e)
		name      = utils.FetchValue(e.Data.Options, "name")
		gID       = generateGladiatorID(mID, name)
		url       = os.Getenv("ARENA_URL")
		a         = arena.NewClient(url)
		data      api.InteractionResponse
	)

	a.GetGladiators(mID)

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
			Content:    option.NewNullableString(msg),
			Components: &[]discord.Component{},
		},
	}

	return data, nil
}

func FireGladiator(e *gateway.InteractionCreateEvent) (api.InteractionResponse, error) {

	var (
		formatMsg = "You have fired %s !"
		msg       string
		mID       = generateManagerID(e)
		name      = utils.FetchValue(e.Data.Options, "name")
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
			Content: option.NewNullableString(msg),
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
