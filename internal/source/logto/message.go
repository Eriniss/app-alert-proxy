package logto

import (
	"fmt"
	"time"

	"github.com/Eriniss/app-alert-proxy/internal/notify"
)

func buildMessage(ev SignInEvent, loc *time.Location) notify.Message {
	app := "알 수 없음"
	if ev.Application != nil && ev.Application.Name != "" {
		app = ev.Application.Name
	}

	return notify.Message{
		Level: notify.LevelInfo,
		Title: "로그인 알림",
		Fields: []notify.Field{
			{Name: "누가", Value: fmt.Sprintf("%s (%s)", ev.User.Name, ev.User.Username)},
			{Name: "언제", Value: ev.CreatedAt.In(loc).Format("2006-01-02 15:04:05 MST")},
			{Name: "어디서", Value: ev.UserIP},
			{Name: "앱", Value: app},
			{Name: "브라우저", Value: ev.UserAgent},
		},
	}
}
