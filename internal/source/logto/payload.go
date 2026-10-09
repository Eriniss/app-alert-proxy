package logto

import (
	"strings"
	"time"
)

type SignInEvent struct {
	Event            string       `json:"event"`
	InteractionEvent string       `json:"interactionEvent"`
	CreatedAt        time.Time    `json:"createdAt"`
	UserAgent        string       `json:"userAgent"`
	UserIP           string       `json:"userIp"`
	User             User         `json:"user"`
	Application      *Application `json:"application"` // 없을 수 있어서 포인터
}

type User struct {
	ID         string         `json:"id"`
	Username   string         `json:"username"`
	Name       string         `json:"name"`
	CustomData map[string]any `json:"customData"`
}

type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (u User) AlertEnabled() bool {
	switch v := u.CustomData["isAlert"].(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true")
	default:
		return false
	}
}
