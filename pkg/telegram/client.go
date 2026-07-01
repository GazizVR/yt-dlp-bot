package telegram

import (
	"fmt"
)

type Client struct {
	Token   string
	BaseURL string
}

func NewClient(
	token string,
	baseURL string,
) *Client {
	return &Client{
		Token:   token,
		BaseURL: baseURL,
	}
}

func (c *Client) urlPath(method string) string {
	return fmt.Sprint("bot", c.Token, "/", method)
}
