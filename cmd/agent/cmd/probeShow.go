package cmd

import (
	"fmt"

	"github.com/eldius/rpi-system-monitor/internal/adapter"
	"github.com/spf13/cobra"
)

// probeShowCmd represents the show command
var probeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display time series data for probe measurements",
	Long:  `Display time series data for probe measurements.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("show called")
		values, err := adapter.Get(cmd.Context())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("######################################################")
		fmt.Println("")
		fmt.Println("Monitoring time series:")
		for _, result := range values {
			fmt.Println("---")
			fmt.Printf("- Timestamp:        %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Printf("  CPU Count:        %d\n", result.CPU.CPUCount)
			fmt.Printf("  CPU Usage:        %.2f%%\n", result.CPU.CPUUsage)
			fmt.Printf("  Memory Usage:     %.2d/%d (%.2f%%)\n", result.Memory.UsedMemory, result.Memory.TotalMemory, result.Memory.MemoryUsagePercentage)
			fmt.Printf("  Memory Usage (h): %s/%s (%.2f%%)\n", result.Memory.UsedMemoryStr(), result.Memory.TotalMemoryStr(), result.Memory.MemoryUsagePercentage)
			fmt.Printf("  Temperature:      %.2f°C\n", result.Temp.Temperature)
		}
		fmt.Println("######################################################")
		fmt.Println("")

	},
}

func init() {
	probeCmd.AddCommand(probeShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeShowCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeShowCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
