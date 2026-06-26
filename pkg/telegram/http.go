package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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
	defer resp.Body.Close()
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

func sendVideoRequest[T any](
	baseURL string,
	urlPath string,
	params map[string]string,
	file os.File,
	v T,
) ([]byte, error) {
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	part, err := writer.CreateFormFile("video", file.Name())
	if err != nil {
		log.Println("Ошибка создания multipart file: ", err)
		return nil, err
	}
	_, err = io.Copy(part, &file)
	if err != nil {
		log.Println("Ошибка записи данных в part: ", err)
		return nil, err
	}
	if err := writer.Close(); err != nil {
		log.Println("Ошибка закрытия multipart writer: ", err)
		return nil, err
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Println("Ошибка парсинга baseURL: ", err)
		return nil, err
	}

	u.Path = urlPath

	urlStr := u.String()

	request, err := http.NewRequest(
		http.MethodPost,
		urlStr,
		body,
	)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.DefaultClient

	resp, err := client.Do(request)
	if err != nil {
		log.Println("Ошибка выполнения post http запроса: ", err)
		return nil, err
	}
	defer resp.Body.Close()
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
	resp CommonResponse,
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
