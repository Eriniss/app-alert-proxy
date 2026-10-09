package notify

import "context"

type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

type Field struct {
	Name  string
	Value string
}

type Message struct {
	Level  Level
	Title  string
	Fields []Field
}

type Notifier interface {
	Send(ctx context.Context, msg Message) error
}
