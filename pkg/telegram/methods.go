package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	getUpdatesMethod      = "getUpdates"
	sendMessageMethod     = "sendMessage"
	editMediaMethod       = "editMessageMedia"
	editReplyMarkupMethod = "editMessageReplyMarkup"
	editTextMethod        = "editMessageText"
	answerQueryMethod     = "answerCallbackQuery"
	sendAudioMethod       = "sendAudio"
)

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
		c.urlPath(getUpdatesMethod),
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
		getUpdatesMethod,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SendMessage(
	chatId int64,
	text string,
	markup *InlineMarkup,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
		"text":    text,
	}
	if markup != nil {
		jsonMarkup, err := json.Marshal(markup)
		if err != nil {
			log.Println("Ошибка преобразования markup json: ", err)
			return nil, err
		}
		params["reply_markup"] = string(jsonMarkup)
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(sendMessageMethod),
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
		sendMessageMethod,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) EditMessageMedia(
	chatId int64,
	messageId int64,
	video os.File,
	markup *InlineMarkup,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id":    fmt.Sprintf("%d", chatId),
		"message_id": fmt.Sprintf("%d", messageId),
		"media": `{
			"type": "video",
			"media": "attach://video"
		}`,
	}
	if markup != nil {
		jsonMarkup, err := json.Marshal(markup)
		if err != nil {
			log.Println("Ошибка преобразования markup json: ", err)
			return nil, err
		}
		params["reply_markup"] = string(jsonMarkup)
	}

	body, err := postFormRequest(
		c.BaseURL,
		c.urlPath(editMediaMethod),
		params,
		map[string]os.File{"video": video},
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		editMediaMethod,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) EditMessageReplyMarkup(
	chatId int64,
	messageId int64,
	markup InlineMarkup,
) (*MessageResponse, error) {
	var response MessageResponse
	jsonMarkup, err := json.Marshal(markup)
	if err != nil {
		log.Println("Ошибка преобразования markup json: ", err)
		return nil, err
	}
	params := map[string]string{
		"chat_id":      fmt.Sprintf("%d", chatId),
		"message_id":   fmt.Sprintf("%d", messageId),
		"reply_markup": string(jsonMarkup),
	}

	body, err := getRequest(
		c.BaseURL,
		c.urlPath(editReplyMarkupMethod),
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
		editReplyMarkupMethod,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SendAudio(
	chatId int64,
	audio os.File,
	messageIdToReply *int64,
) (*MessageResponse, error) {
	var response MessageResponse
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
	}
	if messageIdToReply != nil {
		params["reply_parameters"] = fmt.Sprintf(
			`{"message_id": %d}`,
			*messageIdToReply,
		)
	}

	body, err := postFormRequest(
		c.BaseURL,
		c.urlPath(sendAudioMethod),
		params,
		map[string]os.File{"audio": audio},
		&response,
	)
	if err != nil {
		return nil, err
	}

	resp := CommonResponse{Ok: response.Ok}
	if err := checkError(
		resp,
		body,
		sendAudioMethod,
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
		c.urlPath(editTextMethod),
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
		editTextMethod,
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
		c.urlPath(answerQueryMethod),
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
		answerQueryMethod,
	); err != nil {
		return nil, err
	}
	return &response, nil
}
