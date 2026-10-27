package bot

import (
	"context"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

type Route struct {
	Match  func(telegrambot.Update) bool
	Handle func(ctx context.Context, upd telegrambot.Update) error
}

type Router struct {
	routes []Route
}

func NewRouter(routes []Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) Route(
	ctx context.Context,
	update telegrambot.Update,
) error {
	for _, r := range r.routes {
		if r.Match(update) {
			if err := r.Handle(ctx, update); err != nil {
				return err
			}
		}
	}
	return nil
}
