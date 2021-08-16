package client

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/levigross/grequests"
	"github.com/wihrt/idle_arena/manager"
	"github.com/wihrt/idle_arena/utils"
	"go.uber.org/zap"
)

func (c *Client) RegisterManager(mID string, name string, guildID discord.GuildID, difficulty int) (*manager.Manager, error) {

	var (
		url     = []string{c.URL, utils.APIBase, "managers"}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Register manager",
		zap.String("ManagerID", mID),
		zap.Int("Difficulty", difficulty),
		zap.String("URL", fullURL),
	)

	m, err := manager.NewManager(mID, name, guildID, difficulty)
	if err != nil {
		zap.L().Error("Cannot create manager",
			zap.String("ManagerID", mID),
			zap.Int("Difficulty", difficulty),
			zap.Error(err),
		)
	}

	res, err := grequests.Post(fullURL, &grequests.RequestOptions{
		JSON:           m,
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Error("Cannot register manager",
			zap.String("ManagerID", mID),
			zap.Int("Difficulty", difficulty),
			zap.String("URL", fullURL),
			zap.Error(err),
		)
		return m, err
	}
	if !res.Ok {
		zap.L().Error("Cannot register manager",
			zap.String("ManagerID", mID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return m, ErrWrongStatusCode
	}
	return m, nil
}

func (c *Client) ShowManager(mID string) (*manager.Manager, error) {
	var (
		m       manager.Manager
		url     = []string{c.URL, utils.APIBase, "managers", mID}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Get manager",
		zap.String("ManagerID", mID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Get(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Error("Cannot get manager",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return &m, err
	}
	if !res.Ok {
		zap.L().Error("Cannot retire manager",
			zap.String("ManagerID", mID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return &m, err
	}

	err = res.JSON(&m)
	if err != nil {
		zap.L().Error("Cannot decode document",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return &m, err
	}

	return &m, nil
}

func (c *Client) RetireManager(mID string) error {

	var (
		url     = []string{c.URL, utils.APIBase, "managers", mID}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Retire manager",
		zap.String("ManagerID", mID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Delete(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Error("Cannot retire  manager",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return err
	}
	if !res.Ok {
		zap.L().Error("Cannot retire manager",
			zap.String("ManagerID", mID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return err
	}
	return nil
}