package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /webhook/logto", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			logger.Error("read body failed", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		logger.Info("logto webhook", "headers", r.Header, "body", string(body))
		w.WriteHeader(http.StatusOK)
	})

	addr := ":8080"
	logger.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
