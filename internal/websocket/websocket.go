package websocket

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/dermaddis/op_tournament/internal/errs"
	"github.com/dermaddis/op_tournament/internal/handler/customcontext"
	"github.com/dermaddis/op_tournament/internal/websocket/bus"
	"github.com/dermaddis/op_tournament/internal/websocket/connection"
	"github.com/dermaddis/op_tournament/internal/websocket/message"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebsocketHandler struct {
	outgoingBus *bus.Bus
	log         *slog.Logger

	userConnections   map[string]struct{}
	userConnectionsMu *sync.Mutex

	messages chan message.IncomingMessage

	stopCtx  context.Context
	stopFunc context.CancelFunc

	upgrader *websocket.Upgrader
}

func New(outgoingBus *bus.Bus, log *slog.Logger) *WebsocketHandler {
	ctx, cancel := context.WithCancel(context.Background())

	return &WebsocketHandler{
		outgoingBus: outgoingBus,
		log:         log,

		userConnections:   map[string]struct{}{},
		userConnectionsMu: &sync.Mutex{},

		messages: make(chan message.IncomingMessage),

		stopCtx:  ctx,
		stopFunc: cancel,

		upgrader: &websocket.Upgrader{},
	}
}

func (h *WebsocketHandler) Serve() {
listen:
	for {
		select {
		case <-h.stopCtx.Done():
			break listen
		case message, ok := <-h.messages:
			if !ok {
				break listen
			}

			go h.handleMessage(message)
		}
	}
}

func (h *WebsocketHandler) Stop() error {
	h.stopFunc()

	time.Sleep(100 * time.Millisecond)

	return nil
}

func (h *WebsocketHandler) handleMessage(message message.IncomingMessage) {
	payloadAny, err := message.Parse()
	if err != nil {
		h.log.Error("failed to parse message", slog.Any("err", err))
		return
	}

	switch payload := payloadAny.(type) {
	default:
		h.log.Debug("payload", slog.Any("payload", payload))
	}
	if err != nil {
		var customErr errs.CustomError
		if !errors.As(err, &customErr) {
			h.log.Error("", slog.Any("err", err))
		}
	}
}

func (h *WebsocketHandler) NewConnection(c echo.Context) error {
	customContext, ok := c.(*customcontext.CustomContext)
	if !ok {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	discordId := string(customContext.DiscordUser.Id)

	h.userConnectionsMu.Lock()
	if _, ok := h.userConnections[discordId]; ok {
		h.userConnectionsMu.Unlock()
		return c.String(http.StatusTooManyRequests, "only one connection at a time")
	}
	h.userConnections[discordId] = struct{}{}
	h.userConnectionsMu.Unlock()

	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	ws.SetReadDeadline(time.Time{})
	ws.SetWriteDeadline(time.Time{})

	conn := connection.New(ws, discordId, h.log)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		conn.ReadIncoming(h.stopCtx, h.messages)
		h.log.Debug("connection closed", slog.String("discordId", discordId))

		cancel()

		h.userConnectionsMu.Lock()
		delete(h.userConnections, discordId)
		h.userConnectionsMu.Unlock()
	}()

	outgoing := h.outgoingBus.Subscribe(ctx, discordId)
	go conn.WriteOutging(ctx, outgoing)
	h.log.Debug("new connection", slog.String("discordId", discordId))

	return nil
}
