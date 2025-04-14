package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dermaddis/op_tournament/internal/handler"
)

func Run() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := handler.New(log)

	shutdownCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

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
}
