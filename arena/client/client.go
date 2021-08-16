package client

import (
	"errors"
)

var ErrWrongStatusCode = errors.New("wrong status code")

type Client struct {
	URL string
}

func NewClient(url string) *Client {
	a := &Client{
		URL: url,
	}

	return a
}
