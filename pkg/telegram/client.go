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
	GetUpdates             = "getUpdates"
	SendMessage            = "sendMessage"
	EditMessageMedia       = "editMessageMedia"
	EditMessageReplyMarkup = "editMessageReplyMarkup"
	EditMessageText        = "editMessageText"
	AnswerCallBackQuery    = "answerCallbackQuery"
	SendAudio              = "sendAudio"
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
		c.urlPath(GetUpdates),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		GetUpdates,
	); err != nil {
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
		c.urlPath(SendMessage),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		SendMessage,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) EditMessageToVideo(
	chatId int64,
	messageId int64,
	video os.File,
	btnText string,
	btnCallback string,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id":    fmt.Sprintf("%d", chatId),
		"message_id": fmt.Sprintf("%d", messageId),
		"media": `{
			"type": "video",
			"media": "attach://video"
		}`,
		"reply_markup": fmt.Sprintf(
			`{
				"inline_keyboard": [
					[
						{"text": "%s", "callback_data": "%s"}
					]
				]
			}`,
			btnText,
			btnCallback,
		),
	}

	body, err := postRequest(
		c.BaseURL,
		c.urlPath(EditMessageMedia),
		params,
		video,
		"video",
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		EditMessageMedia,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) DeleteVideoKeyboard(
	chatId int64,
	messageId int64,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id":      fmt.Sprintf("%d", chatId),
		"message_id":   fmt.Sprintf("%d", messageId),
		"reply_markup": `{"inline_keyboard": [[]]}`,
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(EditMessageReplyMarkup),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		EditMessageReplyMarkup,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SendAudio(
	chatId int64,
	audio os.File,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
	}

	body, err := postRequest(
		c.BaseURL,
		c.urlPath(SendAudio),
		params,
		audio,
		"audio",
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		SendAudio,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) EditMessageText(
	chatId int64,
	messageId int64,
	text string,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"message_id": fmt.Sprintf("%d", messageId),
		"chat_id":    fmt.Sprintf("%d", chatId),
		"text":       text,
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(EditMessageText),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		EditMessageText,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) AnserCallbackQuery(
	queryId string,
) (*CommonResponse, error) {
	var response CommonResponse
	params := map[string]string{
		"callback_query_id": queryId,
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(AnswerCallBackQuery),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		AnswerCallBackQuery,
	); err != nil {
		return nil, err
	}
	return &response, nil
}
