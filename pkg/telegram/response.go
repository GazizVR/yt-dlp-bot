package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func getRequest[T any](
	baseURL string,
	urlPath string,
	params map[string]string,
	v T,
) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Println("Ошибка парсинга baseURL: ", err)
		return nil, err
	}

	u.Path = urlPath

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	urlStr := u.String()

	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("Ошибка выполнения get http запроса: ", err)
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка чтения http body: ", err)
		return nil, err
	}
	if err := json.Unmarshal(respBody, v); err != nil {
		log.Println("Ошибка преобразования http body to json: ", err)
		return nil, err
	}

	return respBody, nil
}

func checkError(
	resp CommonResp,
	body []byte,
) error {
	if !resp.Ok {
		var errRespp ErrorResponse
		if err := json.Unmarshal(body, &errRespp); err != nil {
			log.Println("Ошибка преобразования raw bytes to json: ", err)
			return err
		}
		errStr := fmt.Sprintf("Error %d %s\n", errRespp.Code, errRespp.Description)
		err := errors.New(errStr)
		log.Print(err.Error())
		return err
	}
	return nil
}
