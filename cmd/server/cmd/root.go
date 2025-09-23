package cmd

import (
	"github.com/eldius/initial-config-go/configs"
	"github.com/eldius/initial-config-go/setup"
	"github.com/eldius/rpi-system-monitor/internal/config"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpi-telemetry-monitor",
	Short: "A simple tool to monitor the Raspberry Pi system",
	Long:  `A simple tool to monitor the Raspberry Pi system.`,
	PersistentPreRunE: setup.PersistentPreRunE(
		config.AppName,
		setup.WithEnvPrefix("monitor"),
		setup.WithDefaultCfgFileName("config"),
		setup.WithDefaultCfgFileLocations(config.CfgFileLocations...),
		setup.WithConfigFileToBeUsed(cfgFile),
		setup.WithProps(
			config.TemperatureProbeEnabledProp,
		),
		setup.WithDefaultValues(map[string]any{
			configs.LogFormatKey:     configs.LogFormatJSON,
			configs.LogLevelKey:      configs.LogLevelDEBUG,
			configs.LogOutputFileKey: "execution.log",
		}),
	),
	PersistentPostRunE: setup.PersistentPostRunE,
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rpi-telemetry-monitor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
