package service

import (
	"bot/pkg/telegram"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	StartText  = "🔗 Отправьте ссылку на видео"
	SendText   = "⏳ Подождите, загружаем..."
	ErrorText  = "❌ Ошибка установки, попробуйте снова"
	ButtonText = "​📥 Скачать аудио"
)

func (s *Service) handleStartCommand(
	chatId int64,
) error {
	if _, err := s.Tg.SendMessage(
		chatId,
		StartText,
	); err != nil {
		return err
	}
	return nil
}

func (s *Service) handleMsgWURL(
	chatId int64,
	url string,
) error {
	msg, err := s.Tg.SendMessage(
		chatId,
		SendText,
	)
	if err != nil {
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	videoFile, err := s.Dlp.DownloadVideo(
		"tmp",
		url,
	)
	if err != nil {
		s.Tg.DeleteMessage(chatId, msg.Result.Id)
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	s.Tg.DeleteMessage(chatId, msg.Result.Id)
	_, err = s.Tg.SendVideoWithButton(
		chatId,
		*videoFile,
		ButtonText,
		fmt.Sprintf("%s-%d", url, chatId),
	)
	if err != nil {
		s.Tg.DeleteMessage(chatId, msg.Result.Id)
		s.Tg.SendMessage(
			chatId,
			ErrorText,
		)
		return err
	}
	return nil
}

func (s *Service) handleCallbackQuery(
	callback telegram.CallbackQuery,
) error {
	data := callback.Data
	url := data[:strings.LastIndex(data, "-")]
	rawChatId := data[strings.Index(data, "-")+1:]
	chatId, err := strconv.ParseInt(rawChatId, 10, 64)
	if err != nil {
		log.Println("Ошибка конвертации chatId: ", err)
		return err
	}
	s.Tg.DeleteVideoKeyboard(
		callback.Message.Chat.Id,
		callback.Message.Id,
	)
	audio, err := s.Dlp.DownloadAudio("tmp", url)
	if err != nil {
		return err
	}
	s.Tg.SendAudio(chatId, *audio)
	return nil
}
