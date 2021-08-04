package gladiator

import (
	"strings"
	"time"

	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

const URL = "http://names.drycodes.com/1"

func NewRandomName() (string, error) {

	var (
		names []string
		name  string
	)

	res, err := grequests.Get(URL, &grequests.RequestOptions{
		Params: map[string]string{
			"nameOptions": "funnyWords",
			"format":      "json",
		},
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Error("Error when requesting URL",
			zap.String("url", URL),
			zap.Error(err),
		)
		return name, err
	}

	err = res.JSON(&names)
	if err != nil {
		zap.L().Fatal("Cannot decode JSON",
			zap.Error(err),
		)
		return name, err
	}
	name = strings.Replace(names[0], "_", " ", -1)
	return name, nil
}
