package cmd

import (
	"github.com/eldius/initial-config-go/configs"
	"github.com/eldius/initial-config-go/setup"
	"github.com/eldius/rpi-system-monitor/internal/config"
	"github.com/eldius/rpi-system-monitor/internal/system"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpi-system-monitor",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		s := system.NewTelemetry()
		s.Measure(cmd.Context())
	},
	PersistentPreRunE: setup.PersistentPreRunE(
		config.AppName,
		setup.WithEnvPrefix("monitor"),
		setup.WithDefaultCfgFileName("config"),
		setup.WithDefaultCfgFileLocations(config.ConfigFileLocations...),
		setup.WithProps(
			config.TemperatureProbeEnabledProp,
		),
		setup.WithDefaultValues(map[string]any{
			configs.LogFormatKey:     configs.LogFormatJSON,
			configs.LogLevelKey:      configs.LogLevelDEBUG,
			configs.LogOutputFileKey: "execution.log",
		}),
	),
}

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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rpi-system-monitor.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
