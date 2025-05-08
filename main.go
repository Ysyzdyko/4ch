package main

import (
	h "1337b04rd/internal/adapters/left/httptransport"
	"1337b04rd/internal/adapters/right/db"
	"1337b04rd/internal/adapters/right/ricky"
	"1337b04rd/internal/adapters/right/triple_s"
	"1337b04rd/internal/application/api"
	"1337b04rd/internal/application/core"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	stg := db.NewPostgres()
	cfg := triple_s.DefaultConfig()

	smg, err := triple_s.NewTripleSClient(cfg)
	if err != nil {
		slog.Error("failed to init Minio", slog.Any("error", err))
		os.Exit(1)
	}

	ryk := ricky.NewRicky()
	userSvc := core.NewUser() // ✅ создаём реализацию интерфейса UserService

	svc := api.NewApp(smg, userSvc, stg, ryk)

	slog.Info("starting server")
	srv := h.NewHTTPServer(svc) // ✅ создаём сервер
	if err := srv.Serve(); err != nil {
		slog.Error("error serving", slog.Any("error", err))
		os.Exit(1)
	}
}
