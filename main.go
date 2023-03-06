package main

import (
	"fmt"
	"os"

	"github.com/mhbardsley/auto-timetable/backend"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "auto-timetable",
		Short: "Time management program",
		Long:  `Time management program that allows you to generate a timetbale.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: auto-timetable <subcommand> [flags]")
			fmt.Println("subcommands:")
			fmt.Println("  generate - generate a timetable")
			fmt.Println("  help - display this help")
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a timetable",
		Long:  `Generate a timetable from input data`,
		Run: func(cmd *cobra.Command, args []string) {
			fileName, _ := cmd.Flags().GetString("file")
			noOfSlots, _ := cmd.Flags().GetInt("slots")
			threshold, _ := cmd.Flags().GetFloat64("threshold")
			inputData := backend.GetInput(&fileName, noOfSlots)
			backend.GenerateTimetable(inputData, threshold)
		},
	}

	generateCmd.Flags().StringP("file", "f", "input.json", "The input's filename")
	generateCmd.Flags().IntP("slots", "s", 48, "The number of slots to display")
	generateCmd.Flags().Float64P("threshold", "r", 0.04, "Repopulation threshold")

	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
