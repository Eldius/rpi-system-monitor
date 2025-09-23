package telemetry

import (
	"context"
	"github.com/eldius/rpi-system-monitor/internal/model"
	"github.com/shirou/gopsutil/v3/mem"
	"log/slog"
)

func measureMemory(ctx context.Context) model.MemoryResult {
	var result model.MemoryResult
	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		slog.With("error", err).ErrorContext(ctx, "failed to get memory usage")
		return result
	}
	result.MemoryUsagePercentage = vmem.UsedPercent
	result.UsedMemory = int64(vmem.Used)
	result.TotalMemory = int64(vmem.Total)
	return result
}
