package api

import (
	ports "1337b04rd/internal/ports/right"
	"context"
	"sync"
	"time"
)

type App struct {
	sync.Mutex
	userSvc     UserService
	db          ports.DbPort
	minio       ports.MinioPort
	ricky       ports.RickyPort
	timers      map[string]*time.Timer
	ArchivePost func(ctx context.Context, postID string)
}

func NewApp(minio ports.MinioPort, userSvc UserService, db ports.DbPort, ricky ports.RickyPort) *App {
	return &App{
		userSvc: userSvc,
		db:      db,
		ricky:   ricky,
		minio:   minio,
		timers:  make(map[string]*time.Timer),
	}
}

func (a *App) Timers() map[string]*time.Timer {
	return a.timers
}
