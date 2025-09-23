package model

import (
	"fmt"
	"time"
)

type CPUResult struct {
	CPUUsage float64 `json:"cpu_usage"`
	CPUCount int64   `json:"cpu_count"`
}

type MemoryResult struct {
	MemoryUsagePercentage float64 `json:"memory_usage_percentage"`
	UsedMemory            int64   `json:"used_memory"`
	TotalMemory           int64   `json:"total_memory"`
}

func (m MemoryResult) UsedMemoryStr() string {
	return ByteCountIEC(m.UsedMemory)
}

func (m MemoryResult) TotalMemoryStr() string {
	return ByteCountIEC(m.TotalMemory)
}

func (m MemoryResult) MemoryUsagePercentageStr() string {
	return fmt.Sprintf("%.2f%%", m.MemoryUsagePercentage)
}

type TemperatureResult struct {
	Temperature    float64 `json:"temperature"`
	RawTemperature int64   `json:"raw_temperature"`
}

type ProbesResult struct {
	CPU       CPUResult
	Memory    MemoryResult
	Temp      TemperatureResult
	Timestamp time.Time
}

// ByteCountIEC converts a byte count to a human-readable string using IEC (binary) units (base 1024).
func ByteCountIEC(b int64) string {
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
