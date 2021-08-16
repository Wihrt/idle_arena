package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/wihrt/idle_arena/fight"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/manager"
)

func GenerateManagerID(e *gateway.InteractionCreateEvent) string {
	id := manager.GenerateID(e.GuildID.String(), e.Member.User.ID.String())
	return id
}

func GenerateGladiatorID(managerID string, gladiatorName string) string {
	id := gladiator.GenerateID(managerID, gladiatorName)
	return id
}

func ComponentsWrapper(c []discord.Component) []discord.Component {

	var (
		ar         discord.ActionRowComponent
		components []discord.Component
	)

	ar.Components = append(ar.Components, c...)
	components = append(components, ar)
	return components
}

func Button(label string, customID string, style discord.ButtonStyle) discord.ButtonComponent {
	b := discord.ButtonComponent{
		Label:    label,
		CustomID: customID,
		Style:    style,
	}
	return b
}

func GladiatorSelectMenu(g []gladiator.Gladiator, name string, maxValue int) discord.SelectComponent {
	var (
		menuOptions []discord.SelectComponentOption
	)

	if len(g) <= maxValue {
		maxValue = len(g)
	}

	for _, v := range g {
		o := discord.SelectComponentOption{
			Label: v.Name,
			Value: v.Name,
		}
		menuOptions = append(menuOptions, o)
	}

	menu := discord.SelectComponent{
		CustomID:    name,
		Disabled:    false,
		MinValues:   1,
		MaxValues:   maxValue,
		Placeholder: "Choose your gladiator(s)",
		Options:     menuOptions,
	}

	return menu
}

func GladiatorToEmbed(g gladiator.Gladiator) discord.Embed {

	embed := discord.Embed{
		Title: g.Name,
		Type:  discord.NormalEmbed,
		Fields: []discord.EmbedField{
			{Name: g.Strength.Name, Value: strconv.Itoa(g.Strength.Value), Inline: true},
			{Name: g.Dexterity.Name, Value: strconv.Itoa(g.Dexterity.Value), Inline: true},
			{Name: g.Constitution.Name, Value: strconv.Itoa(g.Constitution.Value), Inline: true},
			{Name: "Weapon", Value: g.Weapon.Name, Inline: true},
			{Name: "Armor", Value: g.Armor.Name, Inline: true},
			{Name: "Armor Class", Value: strconv.Itoa(g.ArmorClass), Inline: true},
			{Name: "Level", Value: strconv.Itoa(g.Experience.Level), Inline: true},
			{Name: "Experience", Value: strconv.Itoa(g.Experience.Current) + "/" + strconv.Itoa(g.Experience.NextLevel), Inline: true},
			{Name: "Death Saves", Value: strconv.Itoa(g.DeathSave.Current) + "/" + strconv.Itoa(g.DeathSave.Max), Inline: true},
		},
	}

	return embed
}

func FightToEmbed(f fight.Result) discord.Embed {

	var (
		fightDesc    = "%s has %s the fight !"
		killedDesc   = "%s has lost and has been killed !"
		desc         string
		color        discord.Color
		thumbnailUrl string
	)

	if f.FightWon {
		desc = fmt.Sprintf(fightDesc, f.Gladiator.Name, "won")
		color = 0x00ff00
		thumbnailUrl = "https://www.iconsdb.com/icons/download/green/thumbs-up-24.png"

	} else {
		desc = fmt.Sprintf(fightDesc, f.Gladiator.Name, "lost")
		color = 0xff0000
		thumbnailUrl = "https://www.iconsdb.com/icons/download/red/thumbs-down-24.png"
		if f.KilledInCombat {
			desc = fmt.Sprintf(killedDesc, f.Gladiator.Name)
		}
	}

	embed := discord.Embed{
		Title:       "Fight Result",
		Type:        discord.NormalEmbed,
		Description: desc,
		Color:       color,
		Thumbnail: &discord.EmbedThumbnail{
			URL: thumbnailUrl,
		},
		Fields: []discord.EmbedField{
			{Name: "Money gained", Value: strconv.Itoa(f.MoneyGained), Inline: false},
			{Name: "Enemy level", Value: strconv.Itoa(f.Enemy.Experience.Level), Inline: false},
			{Name: "Enemy " + f.Enemy.Strength.Name, Value: strconv.Itoa(f.Enemy.Strength.Value), Inline: true},
			{Name: "Enemy " + f.Enemy.Dexterity.Name, Value: strconv.Itoa(f.Enemy.Dexterity.Value), Inline: true},
			{Name: "Enemy " + f.Enemy.Constitution.Name, Value: strconv.Itoa(f.Enemy.Constitution.Value), Inline: true},
			{Name: "Enemy Weapon", Value: f.Enemy.Weapon.Name, Inline: true},
			{Name: "Enemy Armor", Value: f.Enemy.Armor.Name, Inline: true},
			{Name: "Enemy Armor Class", Value: strconv.Itoa(f.Enemy.ArmorClass), Inline: true},
		},
	}

	return embed
}

func ManagerToEmbed(m manager.Manager, g []gladiator.Gladiator) discord.Embed {

	m.MoneyPouch.ConvertPieces()
	fields := []discord.EmbedField{
		{Name: "Copper Pieces", Value: strconv.Itoa(m.MoneyPouch.CopperPieces), Inline: true},
		{Name: "Silver Pieces", Value: strconv.Itoa(m.MoneyPouch.SilverPieces), Inline: true},
		{Name: "Gold Pieces", Value: strconv.Itoa(m.MoneyPouch.GoldPieces), Inline: true},
	}

	if len(g) > 0 {
		var gNames []string
		for _, v := range g {
			gNames = append(gNames, v.Name)
		}
		gStr := strings.Join(gNames, "\n")
		gField := discord.EmbedField{Name: "Gladiators", Value: gStr, Inline: false}
		fields = append(fields, gField)
	}

	embed := discord.Embed{
		Title:  "Manager " + m.Name,
		Type:   discord.NormalEmbed,
		Fields: fields,
	}

	return embed
}
