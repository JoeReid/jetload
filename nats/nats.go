package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JoeReid/jetload/specfile"
	"github.com/nats-io/nats.go"
)

var (
	ErrStreamNotFound = errors.New("stream not found")
	ErrStreamNotEmpty = errors.New("stream not empty")
)

type Loader struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func (l *Loader) Close() {
	l.conn.Close()
}

func (l *Loader) Load(ctx context.Context, spec specfile.File, wait bool) error {
	if err := l.checkStream(ctx, spec.Stream); err != nil {
		return err
	}

	if err := l.publishMessages(ctx, spec.Stream, spec.Messages); err != nil {
		return err
	}

	if wait {
		return l.waitForConsumers(ctx, spec.Stream, len(spec.Messages))
	}

	return nil
}

func (l *Loader) checkStream(ctx context.Context, stream string) error {
	var found bool
	for info := range l.js.Streams(nats.Context(ctx)) {
		if info.Config.Name != stream {
			continue
		}

		if info.State.LastSeq != 0 {
			return ErrStreamNotEmpty
		}

		found = true
	}

	if !found {
		return ErrStreamNotFound
	}

	return nil
}

func (l *Loader) waitForConsumers(ctx context.Context, stream string, numMessages int) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(time.Second):
			var messagesOutstanding bool
			for consumer := range l.js.ConsumerNames(stream, nats.Context(ctx)) {
				info, err := l.js.ConsumerInfo(stream, consumer, nats.Context(ctx))
				if err != nil {
					return fmt.Errorf("error getting consumer info for %q: %w", consumer, err)
				}

				if info.NumPending > 0 || info.NumAckPending > 0 {
					fmt.Printf("Consumer %q has %d messages outstanding\n", consumer, int(info.NumPending)+info.NumAckPending)
					messagesOutstanding = true
				}
			}

			if !messagesOutstanding {
				return nil
			}
		}
	}
}

func (l *Loader) publishMessages(ctx context.Context, stream string, messages []specfile.FileMessage) error {
	for i, message := range messages {
		opts := []nats.PubOpt{
			nats.Context(ctx),
			nats.ExpectStream(stream),
			nats.ExpectLastSequence(uint64(i)),
		}

		_, err := l.js.Publish(message.Subject, []byte(message.JSON), opts...)
		if err != nil {
			return fmt.Errorf("error publishing message %d: %w", i, err)
		}
	}

	return nil
}

func NewLoader(natsURL string) (*Loader, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		defer conn.Close()
		return nil, fmt.Errorf("error connecting to jetstream: %w", err)
	}

	return &Loader{
		conn: conn,
		js:   js,
	}, nil
}
