package arena

import (
	"errors"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/wihrt/idle_arena/fight"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/manager"
	"go.uber.org/zap"
)

var ErrWrongStatusCode = errors.New("wrong status code")

type ArenaClient struct {
	URL string
}

func NewArenaClient(url string) *ArenaClient {
	a := &ArenaClient{
		URL: url,
	}

	return a
}

func (a *ArenaClient) RegisterManager(mID string) (*manager.Manager, error) {

	var (
		m       = manager.NewManager(mID)
		url     = []string{a.URL, APIBase, "managers"}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Register manager",
		zap.String("ManagerID", mID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Post(fullURL, &grequests.RequestOptions{
		JSON:           m,
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Error("Cannot register manager",
			zap.String("ManagerID", mID),
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

func (a *ArenaClient) RetireManager(mID string) error {

	var (
		url     = []string{a.URL, APIBase, "managers", mID}
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

func (a *ArenaClient) HireGladiator(mID string) (gladiator.Gladiator, error) {

	var (
		g       gladiator.Gladiator
		url     = []string{a.URL, APIBase, "managers", mID, "gladiators"}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Hiring new gladiator",
		zap.String("ManagerID", mID),
	)

	res, err := grequests.Post(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})

	if err != nil {
		zap.L().Error("Cannot hire new gladiator",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return g, err
	}

	if !res.Ok {
		zap.L().Error("Cannot hire new gladiator",
			zap.String("ManagerID", mID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return g, ErrWrongStatusCode
	}

	err = res.JSON(&g)
	if err != nil {
		zap.L().Error("Cannot hire new gladiator",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return g, err
	}

	zap.L().Debug("Gladiator hired",
		zap.String("GladiatorID", g.GladiatorID),
		zap.String("ManagerID", g.ManagerID),
	)

	return g, nil

}

func (a *ArenaClient) GetGladiators(mID string) ([]gladiator.Gladiator, error) {
	var (
		g       []gladiator.Gladiator
		url     = []string{a.URL, APIBase, "managers", mID, "gladiators"}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Get gladiators",
		zap.String("ManagerID", mID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Get(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})

	if err != nil {
		zap.L().Error("Cannot get gladiators",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return g, err
	}

	if !res.Ok {
		zap.L().Error("Cannot get gladiators",
			zap.String("ManagerID", mID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return g, ErrWrongStatusCode
	}

	err = res.JSON(&g)
	if err != nil {
		zap.L().Error("Cannot get gladiators",
			zap.String("ManagerID", mID),
			zap.Error(err),
		)
		return g, err
	}

	return g, nil

}

func (a *ArenaClient) GetGladiator(mID string, gID string) (gladiator.Gladiator, error) {

	var (
		g       gladiator.Gladiator
		url     = []string{a.URL, APIBase, "managers", mID, "gladiators", gID}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Get gladiator",
		zap.String("ManagerID", mID),
		zap.String("GladiatorID", gID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Get(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})

	if err != nil {
		zap.L().Error("Cannot get gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Error(err),
		)
		return g, err
	}

	if !res.Ok {
		zap.L().Error("Cannot get gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return g, ErrWrongStatusCode
	}

	err = res.JSON(&g)
	if err != nil {
		zap.L().Error("Cannot get gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Error(err),
		)
		return g, err
	}

	return g, nil
}

func (a *ArenaClient) FightGladiator(mID string, gID string) (fight.FightResult, error) {

	var (
		f       fight.FightResult
		url     = []string{a.URL, APIBase, "managers", mID, "gladiators", gID, "fight"}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Fight gladiator",
		zap.String("ManagerID", mID),
		zap.String("GladiatorID", gID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Post(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})

	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Error(err),
		)
		return f, err
	}

	if !res.Ok {
		zap.L().Error("Cannot fight gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return f, ErrWrongStatusCode
	}

	err = res.JSON(&f)
	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Error(err),
		)
		return f, err
	}

	return f, nil
}

func (a *ArenaClient) FireGladiator(mID string, gID string) error {

	var (
		url     = []string{a.URL, APIBase, "managers", mID, "gladiators", gID}
		fullURL = strings.Join(url, "/")
	)

	zap.L().Info("Fire gladiator",
		zap.String("ManagerID", mID),
		zap.String("GladiatorID", gID),
		zap.String("URL", fullURL),
	)

	res, err := grequests.Delete(fullURL, &grequests.RequestOptions{
		RequestTimeout: 5 * time.Second,
	})

	if err != nil {
		zap.L().Error("Cannot fire gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Error(err),
		)
		return err
	}

	if !res.Ok {
		zap.L().Error("Cannot fire gladiator",
			zap.String("ManagerID", mID),
			zap.String("GladiatorID", gID),
			zap.Int("status code", res.StatusCode),
			zap.Error(ErrWrongStatusCode),
		)
		return ErrWrongStatusCode
	}

	return nil
}
