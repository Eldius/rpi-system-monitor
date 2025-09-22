package config

import (
	"github.com/eldius/initial-config-go/setup"
	"github.com/spf13/viper"
)

const (
	AppName = "rpi-monitor-server"
)

var (
	Version   string
	BuildDate string
	Commit    string
)

var (
	TemperatureProbeEnabledProp = setup.Prop{
		Key:   "monitor.temperature.enabled",
		Value: true,
	}

	ConfigFileLocations = []string{
		"~/.config/rpi-monitor",
		"/etc/rpi-monitor",
		".",
		"./config",
	}
)

func GetTemperatureProbeEnabled() bool {
	return viper.GetBool(TemperatureProbeEnabledProp.Key)
}

func GetVersionInfo() map[string]string {
	return map[string]string{
		"Version":   Version,
		"BuildDate": BuildDate,
		"Commit":    Commit,
		"AppName":   AppName,
	}
}
