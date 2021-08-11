package utils

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/gateway"
	"go.uber.org/zap"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringContains(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func FetchValue(options []gateway.InteractionOption, optionName string) string {
	var value string

	for _, o := range options {
		if o.Name == optionName {
			value = o.Value
		}
	}

	zap.L().Debug("Result of fetch value",
		zap.String("name", optionName),
		zap.String("value", value),
	)

	return value
}
