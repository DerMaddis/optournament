package message

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dermaddis/op_tournament/internal/websocket/events"
)

type OutgoingMessage struct {
	ToDiscordIds []string
	Payload      []byte
}

func NewOutgoingBase(toDiscordIds []string, event events.Event) (OutgoingMessage, error) {
	payload := PayloadOut{
		Event: event,
		Data:  nil,
	}
	return NewOutgoing(toDiscordIds, payload)
}

func NewOutgoing(toDiscordIds []string, payload PayloadOut) (OutgoingMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return OutgoingMessage{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return OutgoingMessage{
		ToDiscordIds: toDiscordIds,
		Payload:      payloadBytes,
	}, nil
}

type IncomingMessage struct {
	DiscordId string
	Payload   []byte
}

type PayloadIn struct {
	Event events.Event    `json:"event"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type PayloadOut struct {
	Event events.Event `json:"event"`
	Data  any          `json:"data,omitempty"`
}

var ErrInvalidMessage = errors.New("invalid message")

// Parse parses the message payload based on the message's event
func (m IncomingMessage) Parse() (any, error) {
	var base PayloadIn
	err := json.Unmarshal(m.Payload, &base)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message as BasePayload: %w", err)
	}

	var parseErr error = nil
	var parsed any = nil

	switch base.Event {
	default:
		parseErr = ErrInvalidMessage
		parsed = nil
	}

	return parsed, parseErr
}
