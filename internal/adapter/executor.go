package adapter

import (
	"context"
	"github.com/eldius/rpi-system-monitor/internal/model"
	"github.com/eldius/rpi-system-monitor/internal/persistence"
	"github.com/eldius/rpi-system-monitor/internal/telemetry"
)

func Measure(ctx context.Context) (model.ProbesResult, error) {
	probesResult := telemetry.Measure(ctx)
	return probesResult, persistence.Persist(ctx, &probesResult)
}

func Get(ctx context.Context) ([]model.ProbesResult, error) {
	return persistence.Get(ctx)
}
