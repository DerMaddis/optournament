package connection

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dermaddis/op_tournament/internal/websocket/message"
	"github.com/gorilla/websocket"
)

type Connection struct {
	conn      *websocket.Conn
	DiscordId string
	log       *slog.Logger
}

func New(conn *websocket.Conn, discordId string, log *slog.Logger) *Connection {
	return &Connection{
		conn:      conn,
		DiscordId: discordId,
		log:       log,
	}
}

// ReadIncoming sends any incoming socket message to the channel.
func (c *Connection) ReadIncoming(ctx context.Context, to chan<- message.IncomingMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, payload, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.log.Error("failed to read message", slog.String("discordId", c.DiscordId), slog.Any("err", err))
				}
				return
			}

			// We use json messages. They are still strings at this point.
			if messageType != websocket.TextMessage {
				c.log.Error("wrong message type", slog.String("discordId", c.DiscordId))
				continue
			}

			message := message.IncomingMessage{
				DiscordId: c.DiscordId,
				Payload:   payload,
			}

			to <- message
		}
	}
}

// WriteOutging sends any received message from the channel to the socket.
func (c *Connection) WriteOutging(ctx context.Context, from <-chan []byte) {
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case bytes, ok := <-from:
			// There is the case that the bytes get received and then,
			// over the course of this select case, the connection is closed
			// and the context is cancelled. This results in an expected error
			// in Send. This race condition cannot be fully prevented so we
			// fall back to context-sensitive error handling in Send.
			if !ok {
				return
			}
			err := c.Send(ctx, bytes)
			if err != nil {
				c.log.Error("failed to send message", slog.String("discordId", c.DiscordId), slog.Any("err", err))
			}
		}
	}
}

func (c *Connection) Send(ctx context.Context, message []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	// Context is still open so the error is not related to closure
	if err != nil && ctx.Err() == nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
