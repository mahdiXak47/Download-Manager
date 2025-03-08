package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(queueCmd)

	// Add flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.download-manager.yaml)")
	addCmd.Flags().StringP("url", "u", "", "URL to download")
	addCmd.Flags().StringP("queue", "q", "default", "Queue to add the download to")
	addCmd.Flags().Int64P("speed", "s", 0, "Speed limit in bytes per second (0 for unlimited)")
	addCmd.Flags().StringP("schedule", "t", "", "Schedule time (format: HH:MM or YYYY-MM-DD HH:MM)")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Download Manager v0.1")
	},
}

var addCmd = &cobra.Command{
	Use:   "add [url]",
	Short: "Add a new download",
	Long:  `Add a new download to a specific queue with optional speed limit and schedule.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		queue, _ := cmd.Flags().GetString("queue")
		speed, _ := cmd.Flags().GetInt64("speed")
		schedule, _ := cmd.Flags().GetString("schedule")

		if url == "" && len(args) > 0 {
			url = args[0]
		}

		if url == "" {
			fmt.Println("Error: URL is required")
			return
		}

		fmt.Printf("Adding download: %s to queue: %s (Speed: %d, Schedule: %s)\n", 
			url, queue, speed, schedule)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all downloads",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all downloads...")
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause [download-id]",
	Short: "Pause a download",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: download-id is required")
			return
		}
		fmt.Printf("Pausing download: %s\n", args[0])
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume [download-id]",
	Short: "Resume a download",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: download-id is required")
			return
		}
		fmt.Printf("Resuming download: %s\n", args[0])
	},
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Manage download queues",
	Long:  `Create, list, modify, and delete download queues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Queue management - Feature coming soon!")
	},
} 