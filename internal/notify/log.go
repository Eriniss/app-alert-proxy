package notify

import (
	"context"
	"log/slog"
)

type LogNotifier struct {
	logger *slog.Logger
}

func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{logger: logger}
}

// 컴파일 시점에 Notifier 인터페이스를 만족하는지 확인한다.
var _ Notifier = (*LogNotifier)(nil)

func (n *LogNotifier) Send(ctx context.Context, msg Message) error {
	attrs := []any{"level", msg.Level, "title", msg.Title}
	for _, f := range msg.Fields {
		attrs = append(attrs, f.Name, f.Value)
	}
	n.logger.Info("notification", attrs...)
	return nil
}
