package telemetry

import (
	"context"
	"github.com/eldius/rpi-system-monitor/internal/config"
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
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		result.CPU = measureCPU(ctx)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		result.Memory = measureMemory(ctx)
	}(&wg)

	if config.GetTemperatureProbeEnabled() {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			result.Temp = measureTemperature()
		}(&wg)
	}

	result.Timestamp = time.Now()

	wg.Wait()

	return result
}
