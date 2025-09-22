package system

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	temperatureURI = "/sys/class/thermal/thermal_zone0/temp"
)

type Telemetry struct {
}

func NewTelemetry() *Telemetry {
	return &Telemetry{}
}

func (t *Telemetry) Measure(ctx context.Context) {
	// Get CPU usage
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Printf("Error getting CPU usage: %v\n", err)
		return
	}
	// Get memory usage
	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		fmt.Printf("Error getting memory usage: %v\n", err)
		return
	}

	fmt.Println("")
	fmt.Println("######################################################")

	fmt.Printf("CPU Usage: %.2f%%\n", percentages[0])

	fmt.Printf("Memory Usage: %.2d/%d (%.2f%%)\n", vmem.Used, vmem.Total, vmem.UsedPercent)
	fmt.Printf("Memory Usage (h): %s/%s (%.2f%%)\n", ByteCountIEC(vmem.Used), ByteCountIEC(vmem.Total), vmem.UsedPercent)

	fmt.Printf("Temperature: %.2f°C\n", t.Temperature())

	fmt.Println("######################################################")
	fmt.Println("")

	time.Sleep(5 * time.Second) // Monitor every 5 seconds
}

func (t *Telemetry) Temperature() float64 {
	stat, err := os.Stat(temperatureURI)
	if err != nil {
		slog.With("error", err).Error("failed to stat temperature file")
		return -1
	}

	if stat.IsDir() {
		slog.With("path", temperatureURI).Error("temperature file is a directory")
		return -1
	}

	file, err := os.ReadFile(temperatureURI)
	if err != nil {
		slog.With("error", err).Error("failed to read temperature file")
		return -1
	}

	temperature := strings.Trim(string(file), "\n")
	temp, err := strconv.Atoi(temperature)
	if err != nil {
		slog.With("error", err).Error("failed to convert temperature file to int")
		return -1
	}
	return float64(temp) / 1000.0
}

// ByteCountIEC converts a byte count to a human-readable string using IEC (binary) units (base 1024).
func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
