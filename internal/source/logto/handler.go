package logto

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Eriniss/app-alert-proxy/internal/notify"
)

const (
	maxBodyBytes = 1 << 20          // 1MB
	maxEventAge  = 5 * time.Minute  // 이보다 오래된 요청은 재전송으로 간주
	sendTimeout  = 10 * time.Second // 알림 전송 제한 시간
)

type Handler struct {
	signingKey []byte
	notifier   notify.Notifier
	loc        *time.Location
	logger     *slog.Logger
}

func NewHandler(signingKey string, n notify.Notifier, loc *time.Location, logger *slog.Logger) *Handler {
	return &Handler{
		signingKey: []byte(signingKey),
		notifier:   n,
		loc:        loc,
		logger:     logger,
	}
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	// 1. body 읽기
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	if err != nil {
		h.logger.Error("read body failed", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 2. 서명 검증 (파싱보다 먼저, 원본 바이트 기준)
	if !verifySignature(h.signingKey, body, r.Header.Get("Logto-Signature-Sha-256")) {
		h.logger.Warn("invalid signature", "remote", r.RemoteAddr)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. 파싱
	var ev SignInEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		h.logger.Error("parse payload failed", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 4. 오래된 요청 거부 (서명된 요청을 그대로 다시 보내는 공격 방지)
	if age := time.Since(ev.CreatedAt); age > maxEventAge {
		h.logger.Warn("stale event rejected", "age", age.String(), "remote", r.RemoteAddr)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 5. 알림 대상인지 필터
	if ev.Event != "PostSignIn" || ev.InteractionEvent != "SignIn" {
		h.logger.Info("event ignored", "event", ev.Event, "interaction", ev.InteractionEvent)
		w.WriteHeader(http.StatusOK)
		return
	}
	if !ev.User.AlertEnabled() {
		h.logger.Info("alert disabled for user", "user", ev.User.Username)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 6. 먼저 200을 돌려주고, 전송은 백그라운드에서
	w.WriteHeader(http.StatusOK)
	go h.dispatch(ev)
}

func (h *Handler) dispatch(ev SignInEvent) {
	// r.Context()가 아니라 새 context를 쓴다. 요청 context는
	// 핸들러가 반환되는 순간 취소되어서 전송이 바로 끊긴다.
	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	msg := buildMessage(ev, h.loc)
	if err := h.notifier.Send(ctx, msg); err != nil {
		h.logger.Error("send notification failed", "error", err, "user", ev.User.Username)
	}
}
