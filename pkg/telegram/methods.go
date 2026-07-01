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
	msgIdToReply *int64,
) (*MessageResponse, error) {
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
	if msgIdToReply != nil {
		params["reply_parameters"] = fmt.Sprintf(
			`{"message_id": %d}`,
			*msgIdToReply,
		)
	}

	response, err := c.messageRequest(
		params,
		sendMessageMethod,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) SendAudio(
	chatId int64,
	audio os.File,
	messageIdToReply *int64,
) (*MessageResponse, error) {
	params := map[string]string{
		"chat_id": fmt.Sprintf("%d", chatId),
	}
	if messageIdToReply != nil {
		params["reply_parameters"] = fmt.Sprintf(
			`{"message_id": %d}`,
			*messageIdToReply,
		)
	}

	response, err := c.mediaRequest(
		params,
		sendAudioMethod,
		map[string]os.File{"audio": audio},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) EditMessageMedia(
	chatId int64,
	messageId int64,
	video os.File,
	markup *InlineMarkup,
) (*MessageResponse, error) {
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

	response, err := c.mediaRequest(
		params,
		editMediaMethod,
		map[string]os.File{"video": video},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) EditMessageReplyMarkup(
	chatId int64,
	messageId int64,
	markup InlineMarkup,
) (*MessageResponse, error) {
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

	response, err := c.messageRequest(
		params,
		editReplyMarkupMethod,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) EditMessageText(
	chatId int64,
	messageId int64,
	text string,
) (*MessageResponse, error) {
	params := map[string]string{
		"message_id": fmt.Sprintf("%d", messageId),
		"chat_id":    fmt.Sprintf("%d", chatId),
		"text":       text,
	}

	response, err := c.messageRequest(
		params,
		editTextMethod,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) AnserCallbackQuery(
	queryId string,
) (*CommonResponse, error) {
	params := map[string]string{
		"callback_query_id": queryId,
	}

	response, _, err := c.commonRequest(
		params,
		answerQueryMethod,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
