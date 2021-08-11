package dnd

import (
	"time"

	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

type Client struct {
	BaseURL        string
	RequestOptions *grequests.RequestOptions
}

func NewClient(baseUrl string, timeout time.Duration) *Client {
	c := &Client{
		BaseURL: baseUrl,
		RequestOptions: &grequests.RequestOptions{
			RequestTimeout: timeout,
		},
	}

	return c
}

func (c *Client) GetCategories() (*CategoryIndex, error) {
	var (
		pathUrl = "/api/equipment-categories"
		index   CategoryIndex
	)

	res, err := c.get(pathUrl)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("path", pathUrl),
			zap.Error(err),
		)
		return &index, err
	}

	err = res.JSON(&index)
	if err != nil {
		zap.L().Error("Cannot decode JSON",
			zap.String("path", pathUrl),
			zap.String("body", res.String()),
			zap.Error(err),
		)
	}

	return &index, nil
}

func (c *Client) GetCategory(pathUrl string) (*CategoryList, error) {
	var list CategoryList

	res, err := c.get(pathUrl)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("path", pathUrl),
			zap.Error(err),
		)
		return &list, err
	}

	err = res.JSON(&list)
	if err != nil {
		zap.L().Error("Cannot decode JSON",
			zap.String("path", pathUrl),
			zap.String("body", res.String()),
			zap.Error(err),
		)
		return &list, err
	}

	return &list, nil
}

func (c *Client) GetWeapon(pathUrl string) (*Weapon, error) {
	var weapon Weapon

	res, err := c.get(pathUrl)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("path", pathUrl),
			zap.Error(err),
		)
		return &weapon, err
	}

	err = res.JSON(&weapon)
	if err != nil {
		zap.L().Error("Cannot decode JSON",
			zap.String("path", pathUrl),
			zap.String("body", res.String()),
			zap.Error(err),
		)
		return &weapon, err
	}

	zap.L().Info("Weapon",
		zap.String("name", weapon.Name))

	return &weapon, nil
}

func (c *Client) GetArmor(pathUrl string) (*Armor, error) {
	var armor Armor

	res, err := c.get(pathUrl)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("path", pathUrl),
			zap.Error(err),
		)
		return &armor, err
	}

	err = res.JSON(&armor)
	if err != nil {
		zap.L().Error("Cannot decode JSON",
			zap.String("path", pathUrl),
			zap.String("body", res.String()),
			zap.Error(err),
		)
		return &armor, err
	}

	zap.L().Info("Armor",
		zap.String("name", armor.Name))

	return &armor, err
}

func (c *Client) get(pathUrl string) (*grequests.Response, error) {
	var url = c.BaseURL + pathUrl

	res, err := grequests.Get(url, c.RequestOptions)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("url", url),
			zap.Error(err),
		)
		return res, err
	}

	if !res.Ok {
		zap.L().Error("Cannot get url",
			zap.String("url", url),
			zap.Int("status_code", res.StatusCode),
			zap.Error(err),
		)
		return res, err
	}

	return res, nil
}
