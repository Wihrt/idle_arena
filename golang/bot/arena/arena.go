package arena

import (
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

type ArenaClient struct {
	URL string
}

func NewArenaClient(url string) *ArenaClient {
	a := &ArenaClient{
		URL: url,
	}

	return a
}

func (a *ArenaClient) HireGladiator(e *gateway.InteractionCreateEvent) (Gladiator, error) {

	var (
		g Gladiator
	)

	zap.L().Info("Hiring a new gladiator",
		zap.String("UserId", e.Member.User.ID.String()),
		zap.String("GuildId", e.GuildID.String()),
	)

	res, err := grequests.Post(a.URL, &grequests.RequestOptions{
		Data: map[string]string{
			"user_id":  e.Member.User.ID.String(),
			"guild_id": e.GuildID.String(),
		},
		RequestTimeout: 5,
	})

	if err != nil {
		return g, err
	}

	if !res.Ok {
		return g, err
	}

	err = res.JSON(&g)
	if err != nil {
		return g, err
	}

	return g, nil

}

func (a *ArenaClient) GetGladiator(e *gateway.InteractionCreateEvent) (Gladiator, error) {

	var (
		g Gladiator
	)

	zap.L().Info("Get gladiator",
		zap.String("UserId", e.Member.User.ID.String()),
		zap.String("GuildId", e.GuildID.String()),
	)

	res, err := grequests.Get(a.URL, &grequests.RequestOptions{
		Params: map[string]string{
			"user_id":  e.Member.User.ID.String(),
			"guild_id": e.GuildID.String(),
		},
		RequestTimeout: 5,
	})

	if err != nil {
		return g, err
	}

	if !res.Ok {
		return g, err
	}

	err = res.JSON(&g)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (a *ArenaClient) FightGladiator(e *gateway.InteractionCreateEvent) (Gladiator, error) {

	var (
		g Gladiator
	)

	zap.L().Info("Fight gladiator",
		zap.String("UserId", e.Member.User.ID.String()),
		zap.String("GuildId", e.GuildID.String()),
	)

	res, err := grequests.Post(a.URL+"/fight", &grequests.RequestOptions{
		Params: map[string]string{
			"user_id":  e.Member.User.ID.String(),
			"guild_id": e.GuildID.String(),
		},
		RequestTimeout: 5,
	})

	if err != nil {
		return g, err
	}

	if !res.Ok {
		return g, err
	}

	err = res.JSON(&g)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (a *ArenaClient) FireGladiator(e *gateway.InteractionCreateEvent) error {

	zap.L().Info("Fire gladiator",
		zap.String("UserId", e.Member.User.ID.String()),
		zap.String("GuildId", e.GuildID.String()),
	)

	res, err := grequests.Delete(a.URL+"/fire", &grequests.RequestOptions{
		Params: map[string]string{
			"user_id":  e.Member.User.ID.String(),
			"guild_id": e.GuildID.String(),
		},
		RequestTimeout: 5,
	})

	if err != nil {
		return err
	}

	if !res.Ok {
		return err
	}

	return nil
}
