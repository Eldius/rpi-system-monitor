/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/eldius/rpi-system-monitor/internal/tui/simple_charts"
	"github.com/spf13/cobra"
)

// testingCmd represents the testing command
var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "A simple testing subcommand",
	Long:  `A simple testing subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("testing called")
		if err := simple_charts.Start(cmd.Context()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
