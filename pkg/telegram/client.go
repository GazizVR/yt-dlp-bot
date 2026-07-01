package telegram

import (
	"encoding/json"
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

func (c *Client) urlPath(method string) string {
	return fmt.Sprint("bot", c.Token, "/", method)
}

func (c *Client) commonRequest(
	params map[string]string,
	method string,
) (*CommonResponse, []byte, error) {
	var response CommonResponse
	body, err := getRequest(
		c.BaseURL,
		c.urlPath(method),
		params,
		&response,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := checkError(
		response,
		body,
		method,
	); err != nil {
		return nil, nil, err
	}
	return &response, body, nil
}

func (c *Client) messageRequest(
	params map[string]string,
	method string,
) (*MessageResponse, error) {
	var v MessageResponse
	_, respBody, err := c.commonRequest(params, method)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBody, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) mediaRequest(
	params map[string]string,
	method string,
	media map[string]os.File,
) (*MessageResponse, error) {
	var response MessageResponse
	body, err := postFormRequest(
		c.BaseURL,
		c.urlPath(method),
		params,
		media,
		&response,
	)
	if err != nil {
		return nil, err
	}

	if err := checkError(
		response.CommonResponse,
		body,
		method,
	); err != nil {
		return nil, err
	}
	return &response, nil
}
