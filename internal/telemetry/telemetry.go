package telemetry

import (
	"context"
	"github.com/eldius/rpi-system-monitor/internal/feature_toggle"
	"github.com/eldius/rpi-system-monitor/internal/model"
	"sync"
	"time"
)

const (
	temperatureURI = "/sys/class/thermal/thermal_zone0/temp"
)

func Measure(ctx context.Context) model.ProbesResult {

	var result model.ProbesResult

	var wg sync.WaitGroup
	wg.Go(func() {
		result.CPU = measureCPU(ctx)
	})

	wg.Go(func() {
		result.Memory = measureMemory(ctx)
	})

	_ = feature_toggle.FeatureToggle(ctx, "monitor.server.temperature_probe.enabled", func(ctx context.Context) error {
		wg.Go(func() {
			result.Temp = measureTemperature()
		})
		return nil
	})

	result.Timestamp = time.Now()

	wg.Wait()

	return result
}
