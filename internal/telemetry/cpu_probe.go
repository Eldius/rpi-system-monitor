package telemetry

import (
	"context"
	"log/slog"
	"time"

	"github.com/eldius/rpi-system-monitor/internal/model"
	"github.com/shirou/gopsutil/v3/cpu"
)

func measureCPU(ctx context.Context) model.CPUResult {
	var result model.CPUResult
	cpuCount, err := cpu.Counts(false)
	if err != nil {
		slog.With("error", err).ErrorContext(ctx, "failed to get CPU count")
		return result
	}
	result.CPUCount = int64(cpuCount)
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		slog.With("error", err).ErrorContext(ctx, "failed to get CPU usage")
		return result
	}

	result.CPUUsage = percentages[0]

	return result
}
