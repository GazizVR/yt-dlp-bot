package bot

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
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

func (b *Bot) Start(ctx context.Context) error {
	handler := NewHandler(b.client, b.src, b.filesDir)
	updChan, err := b.client.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout:        45,
		AllowedUpdates: []string{"message", "callback_query"},
	})
	if err != nil {
		return err
	}
	thandler, err := th.NewBotHandler(b.client, updChan)
	if err != nil {
		return err
	}
	thandler.Handle(handler.OnStart, th.CommandEqual("start"))
	thandler.HandleMessage(handler.OnDownload, func(ctx context.Context, u telego.Update) bool {
		var match bool
		if u.Message != nil {
			if u.Message.LinkPreviewOptions != nil {
				match = true
			}
		}
		return match
	})
	thandler.HandleCallbackQuery(handler.OnCallbackQuery)
	return thandler.Start()
}
