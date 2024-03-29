package main

import (
	"fmt"
	"os"

	"github.com/mhbardsley/auto-timetable/backend"
	"github.com/mhbardsley/auto-timetable/cli"
	"github.com/spf13/cobra"
)

func main() {
	// recurseively call command-line functions to build up the CLI
	// in a modular fashion
	rootCmd := makeRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
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

	rootCmd.AddCommand(makeGenerateCommand())
	rootCmd.AddCommand(makeAddCommand())

	return rootCmd
}


func makeGenerateCommand() *cobra.Command {
	var dirName string
	var noOfSlots int
	var threshold float64

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a timetable",
		Long:  `Generate a timetable from input data`,
		Run: func(cmd *cobra.Command, args []string) {
			inputData := backend.GetInput(&dirName, noOfSlots)
			backend.GenerateTimetable(inputData, threshold)
		},
	}

	generateCmd.Flags().StringVarP(&dirName, "dir", "d", "toplevel/", "Toplevel directory")
	generateCmd.Flags().IntVarP(&noOfSlots, "slots", "s", 48, "The number of slots to display")
	generateCmd.Flags().Float64VarP(&threshold, "threshold", "r", 0.04, "Repopulation threshold")

	return generateCmd
}

func makeAddCommand() *cobra.Command {
	var dirName string

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add an event or deadline",
		Long:  `Add an event or deadline to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: auto-timetable add [flags]")
			fmt.Println("subsubcommands:")
			fmt.Println("  event - add an event")
			fmt.Println("  deadline - add a deadline")
			fmt.Println("  help - display this help")

			fmt.Println(dirName)
		},
	}

	addCmd.PersistentFlags().StringVarP(&dirName, "file", "f", "input.json", "The input's filename")

	addCmd.AddCommand(makeAddEventCommand())
	addCmd.AddCommand(makeAddDeadlineCommand())

	return addCmd
}

func makeAddEventCommand() *cobra.Command {
	var eventName, startTimeStr, endTimeStr string

	eventCmd := &cobra.Command{
		Use:   "event",
		Short: "Add an event",
		Long:  `Add an event to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			dirName, _ := cmd.Flags().GetString("file")
			// TODO: send via a struct
			cli.AddEvent(&dirName)
		},
	}

	eventCmd.Flags().StringVarP(&eventName, "name", "n", "", "Name of the event")
	eventCmd.Flags().StringVarP(&startTimeStr, "startTime", "s", "", "Start time")
	eventCmd.Flags().StringVarP(&endTimeStr, "endTime", "e", "", "End time")

	return eventCmd
}

func makeAddDeadlineCommand() *cobra.Command {
	var deadlineName, deadlineStr string
	var minutesRemaining float64

	deadlineCmd := &cobra.Command{
		Use:   "deadline",
		Short: "Add a deadline",
		Long:  `Add a deadline to the existing timetable`,
		Run: func(cmd *cobra.Command, args []string) {
			dirName, _ := cmd.Flags().GetString("file")
			// TODO: send via a struct
			cli.AddDeadline(&dirName)
		},
	}

	deadlineCmd.Flags().StringVarP(&deadlineName, "name", "n", "", "Name of the deadline")
	deadlineCmd.Flags().Float64VarP(&minutesRemaining, "minutesRemaining", "m", 25.0, "Time to complete deadline")
	deadlineCmd.Flags().StringVarP(&deadlineStr, "deadline", "d", "", "Time of the deadline")

	return deadlineCmd
}