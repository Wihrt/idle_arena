package main

import (
	"os"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/wihrt/idle_arena/bot/bot"
	"github.com/wihrt/idle_arena/bot/commands"
	"go.uber.org/zap"
)

// To run, do `APP_ID="APP ID" GUILD_ID="GUILD ID" BOT_TOKEN="TOKEN HERE" go run .`

func main() {

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	appID := discord.AppID(mustSnowflakeEnv("APP_ID"))
	guildID := discord.GuildID(mustSnowflakeEnv("GUILD_ID"))

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		zap.L().Fatal("No $BOT_TOKEN given")
	}

	bot, err := bot.NewArenaBot(token, appID, guildID)
	if err != nil {
		zap.L().Fatal("Cannot create bot",
			zap.Error(err),
		)
	}

	bot.AddIntents([]gateway.Intents{gateway.IntentGuilds, gateway.IntentGuildMessages})
	bot.Session.AddHandler(func(e *gateway.InteractionCreateEvent) {
		data := commands.HandleInteraction(e)
		err := bot.Session.RespondInteraction(e.ID, e.Token, data)
		if err != nil {
			zap.L().Error("Failed to send interaction callback",
				zap.Error(err),
			)
		}
	})

	err = bot.Session.Open()
	if err != nil {
		zap.L().Fatal("Failed to open session",
			zap.Error(err),
		)
	}
	defer bot.Session.Close()

	bot.GetGuildCommands()
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
