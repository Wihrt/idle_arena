package bot

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"go.uber.org/zap"
)

type ArenaBot struct {
	Session *session.Session
	AppID   discord.AppID
	GuildID discord.GuildID
}

func NewArenaBot(token string, appID discord.AppID, guildID discord.GuildID) (*ArenaBot, error) {

	var b = &ArenaBot{}

	zap.L().Info("Creating new Arena Bot")
	s, err := session.New("Bot " + token)
	if err != nil {
		zap.L().Fatal("Cannot start bot",
			zap.Error(err),
		)
		return b, err
	}

	b.Session = s
	b.AppID = appID
	b.GuildID = guildID

	return b, nil
}

func (a *ArenaBot) GetGuildCommands() error {

	commands, err := a.Session.GuildCommands(a.AppID, a.GuildID)
	if err != nil {
		zap.L().Fatal("Failed to get guild commands",
			zap.Error(err),
		)
		return err
	}

	for _, command := range commands {
		zap.L().Info("Existing command",
			zap.String("name", command.Name),
			zap.String("description", command.Description),
		)
	}
	return nil
}

func (a *ArenaBot) AddIntents(intents []gateway.Intents) {
	for _, intent := range intents {
		a.Session.Gateway.AddIntents(intent)
	}
}

func (a *ArenaBot) RegisterCommands(commands []api.CreateCommandData) {

	for _, command := range commands {
		_, err := a.Session.CreateGuildCommand(a.AppID, a.GuildID, command)
		if err != nil {
			zap.L().Fatal("Failed to register command",
				zap.String("name", command.Name),
				zap.String("description", command.Description),
			)
		}
	}
}
