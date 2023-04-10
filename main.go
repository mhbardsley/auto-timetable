package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mhbardsley/auto-timetable/backend"
	"github.com/mhbardsley/auto-timetable/cli"
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
			fmt.Println("  add - add an event or deadline")
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

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add an event or deadline",
		Long:  `Add an event or deadline to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: auto-timetable add [flags]")
			fmt.Println("subsubcommands:")
			fmt.Println("  event - add an event")
			fmt.Println("  deadline - add a deadline")
			fmt.Println("  help - display this help")
		},
	}

	// TODO: replace with StringVarP
	addCmd.PersistentFlags().StringP("file", "f", "input.json", "The input's filename")

	var eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Add an event",
		Long:  `Add an event to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			eventName, _ := cmd.Flags().GetString("name")
			repopulate, _ := cmd.Flags().GetBool("repopulate")
			startTimeStr, _ := cmd.Flags().GetString("startTime")
			endTimeStr, _ := cmd.Flags().GetString("endTime")

			startTime, _ := time.Parse(time.RFC3339, startTimeStr)
			endTime, _ := time.Parse(time.RFC3339, endTimeStr)
			fmt.Print(eventName)
			fmt.Print(repopulate)
			fmt.Print(startTime)
			fmt.Print(endTime)
		},
	}

	eventCmd.Flags().StringP("name", "n", "", "Name of the event")
	eventCmd.Flags().BoolP("repopulate", "r", false, "Event is repopulation")
	eventCmd.Flags().StringP("startTime", "s", "", "Start time")
	eventCmd.Flags().StringP("endTime", "e", "", "End time")
	// TODO: flags in the event command

	var deadlineCmd = &cobra.Command{
		Use:   "deadline",
		Short: "Add a deadline",
		Long:  `Add a deadline to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			fileName, _ := cmd.Flags().GetString("file")
			cli.AddDeadline(&fileName)
		},
	}

	addCmd.AddCommand(eventCmd)
	addCmd.AddCommand(deadlineCmd)

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
