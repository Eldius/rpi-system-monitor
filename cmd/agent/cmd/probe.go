/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/eldius/rpi-system-monitor/internal/adapter"

	"github.com/spf13/cobra"
)

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Fetch the current probe values",
	Long:  `Fetch the current probe values.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := adapter.Measure(cmd.Context())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("")
		fmt.Println("######################################################")

		fmt.Printf("CPU Count: %d\n", result.CPU.CPUCount)
		fmt.Printf("CPU Usage: %.2f%%\n", result.CPU.CPUUsage)
		fmt.Printf("Memory Usage: %.2d/%d (%.2f%%)\n", result.Memory.UsedMemory, result.Memory.TotalMemory, result.Memory.MemoryUsagePercentage)
		fmt.Printf("Memory Usage (h): %s/%s (%.2f%%)\n", result.Memory.UsedMemoryStr(), result.Memory.TotalMemoryStr(), result.Memory.MemoryUsagePercentage)
		fmt.Printf("Temperature: %.2f°C\n", result.Temp.Temperature)

		fmt.Println("######################################################")
		fmt.Println("")

	},
}

func init() {
	rootCmd.AddCommand(probeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// probeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// probeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
