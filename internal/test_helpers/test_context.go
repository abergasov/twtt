package testhelpers

import (
	"context"
	"testing"
	"time"
	"twtt/internal/logger"
)

type TestContainer struct {
	Ctx    context.Context
	Logger logger.AppLogger
}

func GetClean(t *testing.T) *TestContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)

	t.Cleanup(func() {
		cancel()
	})

	appLog := logger.NewAppSLogger("test")
	// repo init

	// service init
	return &TestContainer{
		Ctx:    ctx,
		Logger: appLog,
	}
}
