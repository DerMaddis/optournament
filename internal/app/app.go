package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dermaddis/op_tournament/internal/handler"
	"github.com/dermaddis/op_tournament/internal/model/tournament"
	"github.com/dermaddis/op_tournament/internal/service"
	"github.com/dermaddis/op_tournament/internal/websocket"
	"github.com/dermaddis/op_tournament/internal/websocket/bus"
)

func Run() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	wsBus := bus.New()

	service := service.New(log)

	_, _ = service.NewTournament([]tournament.Song{
		{"a"},
		{"b"},
		{"c"},
		{"d"},
		{"e"},
		{"f"},
		{"g"},
		{"h"},
		{"i"},
		{"j"},
		{"k"},
		{"l"},
		{"m"},
		{"n"},
		{"o"},
		{"p"},
	}, "484735504110190602")

	wsHandler := websocket.New(wsBus, log)
	handler := handler.New(service, wsHandler, log)

	shutdownCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		wsHandler.Serve()
	}()

	go func() {
		err := handler.Serve(fmt.Sprintf(":%d", 3000))
		if err != nil {
			log.Error("failed to start rest handler", slog.Any("err", err))
		}
	}()

	<-shutdownCtx.Done()

	log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := handler.Stop(ctx)
	if err != nil {
		log.Error("failed to stop rest handler", slog.Any("err", err))
	}

	err = wsHandler.Stop()
	if err != nil {
		log.Error("failed to stop websocket handler", slog.Any("err", err))
	}
}
