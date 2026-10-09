package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Eriniss/app-alert-proxy/internal/notify"
	"github.com/Eriniss/app-alert-proxy/internal/source/logto"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	signingKey := os.Getenv("LOGTO_SIGNING_KEY")
	if signingKey == "" {
		logger.Error("LOGTO_SIGNING_KEY is required")
		os.Exit(1)
	}

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		logger.Error("load location failed", "error", err)
		os.Exit(1)
	}

	notifier := notify.NewLogNotifier(logger)
	logtoHandler := logto.NewHandler(signingKey, notifier, loc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook/logto", logtoHandler.Webhook)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	addr := ":8080"
	logger.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
