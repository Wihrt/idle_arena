package commands

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func PingCommand(e *gateway.InteractionCreateEvent) api.InteractionResponse {
	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: "Pong!",
		},
	}

	return data

}
