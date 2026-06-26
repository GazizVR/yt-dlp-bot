package telegram

import (
	"fmt"
	"os"
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

const (
	GetUpdatesMethod  = "getUpdates"
	SendMessageMethod = "sendMessage"
	SendVideoMethod   = "sendVideo"
	DeleteMessage     = "deleteMessage"
)

func (c *Client) urlPath(method string) string {
	return fmt.Sprint("bot", c.Token, "/", method)
}

func (c *Client) GetUpdates(
	offset int64,
	limit uint8,
	timeout uint8,
	allowedUpdates []string,
) (*UpdatesResponse, error) {
	var response UpdatesResponse

	params := map[string]string{
		"offset":          fmt.Sprintf("%d", offset),
		"limit":           fmt.Sprintf("%d", limit),
		"timeout":         fmt.Sprintf("%d", timeout),
		"allowed_updates": fmt.Sprintf("%v", allowedUpdates),
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(GetUpdatesMethod),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(resp, body); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SendMessage(
	chatId int64,
	text string,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
		"text":    text,
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(SendMessageMethod),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(resp, body); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SendVideo(
	chatId int64,
	video os.File,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
	}

	body, err := sendVideoRequest(
		c.BaseURL,
		c.urlPath(SendVideoMethod),
		params,
		video,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(resp, body); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) DeleteMessage(
	chatId int64,
	messageId int64,
) error {
	var response CommonResponse
	params := map[string]string{
		"chat_id":    fmt.Sprintf("%d", chatId),
		"message_id": fmt.Sprintf("%d", messageId),
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(SendMessageMethod),
		params,
		&response,
	)
	if err != nil {
		return err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(resp, body); err != nil {
		return err
	}
	return nil
}
