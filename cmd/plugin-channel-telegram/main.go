package main

import (
	"context"
	"log"
	"os"

	"fluxa/pkg/plugin"
	"github.com/jmjp/plugin-channel-telegram/internal/bot"
)

func main() {
	// O servidor gRPC é iniciado pelo supervisor do Fluxa; o loop mantém o
	// polling isolado e pode ser conectado ao host gRPC na composição seguinte.
	if err := bot.Poll(context.Background(), nil, os.Getenv("TELEGRAM_BOT_TOKEN"), func(context.Context, plugin.InboundMessage) error { return nil }); err != nil {
		log.Print(err)
	}
}
