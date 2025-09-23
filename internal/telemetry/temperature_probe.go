package telemetry

import (
	"github.com/eldius/rpi-system-monitor/internal/model"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func measureTemperature() model.TemperatureResult {
	result := model.TemperatureResult{
		Temperature:    -1,
		RawTemperature: -1,
	}

	stat, err := os.Stat(temperatureURI)
	if err != nil {
		slog.With("error", err).Error("failed to stat temperature file")
		return result
	}

	if stat.IsDir() {
		slog.With("path", temperatureURI).Error("temperature file is a directory")
		return result
	}

	file, err := os.ReadFile(temperatureURI)
	if err != nil {
		slog.With("error", err).Error("failed to read temperature file")
		return result
	}

	temperature := strings.Trim(string(file), "\n")
	temp, err := strconv.Atoi(temperature)
	if err != nil {
		slog.With("error", err).Error("failed to convert temperature file to int")
		return result
	}
	result.RawTemperature = int64(temp)
	result.Temperature = float64(temp) / 1000.0
	return result
}
