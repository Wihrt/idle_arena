package main

import (
	"os"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/wihrt/idle_arena/bot/commands"
	"github.com/wihrt/idle_arena/logging"
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

func init() {

	cfg := logging.GetConfig()
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func main() {

	zap.L().Info("Starting bot")

	appID := discord.AppID(mustSnowflakeEnv("APP_ID"))
	guildID := discord.GuildID(mustSnowflakeEnv("GUILD_ID"))
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		zap.L().Fatal("No $BOT_TOKEN given")
	}

	bot, err := NewArenaBot(token, appID, guildID)
	if err != nil {
		zap.L().Fatal("Cannot create bot",
			zap.Error(err),
		)
	}

	bot.AddIntents([]gateway.Intents{gateway.IntentGuilds, gateway.IntentGuildMessages})
	bot.Session.AddHandler(func(e *gateway.InteractionCreateEvent) {
		data, err := commands.HandleInteraction(e)
		if err != nil {
			zap.L().Error("Failed to handle command",
				zap.Error(err),
			)
			return
		}
		err = bot.Session.RespondInteraction(e.ID, e.Token, data)
		if err != nil {
			zap.L().Error("Failed to send interaction callback",
				zap.Error(err),
			)
			return
		}
	})

	err = bot.Session.Open()
	if err != nil {
		zap.L().Fatal("Failed to open session",
			zap.Error(err),
		)
	}
	defer bot.Session.Close()

	err = bot.GetGuildCommands()
	if err != nil {
		zap.L().Fatal("Failed to get guild commands",
			zap.Error(err),
		)
	}
	bot.RegisterCommands(commands.RegisteredCommands)

	// Block forever.
	select {}
}

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		zap.L().Fatal("Invalid Snowflake",
			zap.String("value", env),
			zap.Error(err),
		)
	}
	return s
}
