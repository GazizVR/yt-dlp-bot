package bot

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
)

type Bot struct {
	client   *telego.Bot
	filesDir string
	src      MediaService
}

func NewBot(baseURL, token, filesDir string, src MediaService) (*Bot, error) {
	var opts []telego.BotOption
	if baseURL != "" {
		opts = append(opts, telego.WithAPIServer(baseURL))
	}
	bot, err := telego.NewBot(token, opts...)
	if err != nil {
		return nil, err
	}
	return &Bot{
		filesDir: filesDir,
		client:   bot,
		src:      src,
	}, nil
}

func (b *Bot) ListenUpdates(ctx context.Context) error {
	handler := NewHandler(b.client, b.src, b.filesDir)
	routes := []Route{
		{
			Match: func(u telego.Update) bool {
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
			Match: func(u telego.Update) bool {
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
			Match: func(u telego.Update) bool {
				return u.CallbackQuery != nil
			},
			Handle: handler.OnCallbackQuery,
		},
	}
	router := NewRouter(routes)
	var lastUpdateID int
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		params := &telego.GetUpdatesParams{
			Offset:         lastUpdateID + 1,
			Limit:          100,
			Timeout:        45,
			AllowedUpdates: []string{"message", "callback_query"},
		}
		updates, err := b.client.GetUpdates(ctx, params)
		if err != nil {
			return err
		}
		for _, u := range updates {
			lastUpdateID = u.UpdateID
			go func(u telego.Update) {
				if err := router.Route(ctx, u); err != nil {
					slog.ErrorContext(ctx, err.Error())
				}
			}(u)
		}
	}
}
