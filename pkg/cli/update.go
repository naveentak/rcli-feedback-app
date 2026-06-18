package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rcli/feedback/internal/feedback"
)

var updateCmd = &cobra.Command{
	Use:   "update [number]",
	Short: "Update a feedback ticket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var number int
		if _, err := fmt.Sscanf(args[0], "%d", &number); err != nil {
			return fmt.Errorf("invalid issue number: %s", args[0])
		}

		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("--status is required")
		}

		client := NewAPIClient()
		ticket, err := client.UpdateStatus(number, feedback.Status(status))
		if err != nil {
			return err
		}

		fmt.Printf("Updated #%d → status: %s\n", ticket.Number, ticket.Status)
		return nil
	},
}

func init() {
	updateCmd.Flags().String("status", "", "New status (triaged, in-progress, done)")
}