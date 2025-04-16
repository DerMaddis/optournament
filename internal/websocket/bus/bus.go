package bus

import (
	"context"
	"slices"
	"sync"

	"github.com/dermaddis/op_tournament/internal/websocket/message"
)

type Bus struct {
	subscribers   map[chan<- message.OutgoingMessage]struct{}
	subscribersMu *sync.Mutex
}

func New() *Bus {
	return &Bus{
		subscribers:   map[chan<- message.OutgoingMessage]struct{}{},
		subscribersMu: &sync.Mutex{},
	}
}

func (b *Bus) Send(ctx context.Context, message message.OutgoingMessage) {
	select {
	case <-ctx.Done():
		return
	default:
		b.subscribersMu.Lock()
		defer b.subscribersMu.Unlock()
		for subscriber := range b.subscribers {
			select {
			// We can return since there is no need to guarantee every single message to be sent
			case <-ctx.Done():
				return
			case subscriber <- message:
				// sent
			}
		}
	}
}

func (b *Bus) Subscribe(ctx context.Context, discordId string) <-chan []byte {
	allMessages := make(chan message.OutgoingMessage, 1)
	b.subscribersMu.Lock()
	b.subscribers[allMessages] = struct{}{}
	b.subscribersMu.Unlock()

	go func() {
		<-ctx.Done()
		b.subscribersMu.Lock()
		delete(b.subscribers, allMessages)
		b.subscribersMu.Unlock()
	}()

	toConnection := make(chan []byte, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-allMessages:
				if !ok {
					return
				}
				if slices.Contains(m.ToDiscordIds, discordId) {
					select {
					case <-ctx.Done():
						return
					case toConnection <- m.Payload:
						// sent
					}
				}
			}
		}
	}()

	return toConnection
}
