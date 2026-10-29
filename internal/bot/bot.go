package bot

import (
	"context"
	"log/slog"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

type Bot struct {
	bot      *telegrambot.Bot
	filesDir string
	src      MediaService
}

func NewBot(baseURL, token, filesDir string, src MediaService) (*Bot, error) {
	bot, err := telegrambot.NewBot(token)
	if err != nil {
		return nil, err
	}
	if baseURL != "" {
		if err := telegrambot.WithBaseURL(baseURL)(bot); err != nil {
			return nil, err
		}
	}
	if filesDir != "" {
		if err := telegrambot.WithFilesDir(filesDir)(bot); err != nil {
			return nil, err
		}
	}
	localBot := &Bot{
		filesDir: filesDir,
		bot:      bot,
		src:      src,
	}
	return localBot, nil
}

func (b *Bot) ListenUpdates(ctx context.Context) error {
	handler := NewHandler(b.bot, b.src)
	routes := []Route{
		{
			Match: func(u telegrambot.Update) bool {
				var match bool
				if u.Message != nil {
					if u.Message.Text == "/start" {
						match = true
					}
				}
				return match
			},
			Handle: handler.OnStart,
		},
		{
			Match: func(u telegrambot.Update) bool {
				var match bool
				if u.Message != nil {
					if u.Message.LinkPreviewOptions != nil {
						match = true
					}
				}
				return match
			},
			Handle: handler.OnDownload,
		},
		{
			Match: func(u telegrambot.Update) bool {
				return u.CallbackQuery != nil
			},
			Handle: handler.OnCallbackQuery,
		},
	}
	router := NewRouter(routes)
	var lastUpdateID int64
	for {
		params := telegrambot.GetUpdatesParams{
			Offset:         lastUpdateID + 1,
			Limit:          100,
			Timeout:        45,
			AllowedUpdates: []string{"message", "callback_query"},
		}
		updates, err := b.bot.GetUpdates(ctx, params)
		if err != nil {
			return err
		}
		for _, u := range updates {
			lastUpdateID = u.ID
			go func() {
				if err := router.Route(ctx, u); err != nil {
					slog.ErrorContext(ctx, err.Error())
				}
			}()
		}
	}
}
