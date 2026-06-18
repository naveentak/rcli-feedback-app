package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rcli/feedback/internal/feedback"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a new feedback ticket",
	Example: `  rcli feedback submit --app rclip --type bug --title "Crash on export" --description "Steps to reproduce..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, _ := cmd.Flags().GetString("app")
		ticketType, _ := cmd.Flags().GetString("type")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		reporter, _ := cmd.Flags().GetString("reporter")

		client := NewAPIClient()
		client.app = app

		ticket, err := client.Submit(feedback.SubmitRequest{
			App:         feedback.App(app),
			Type:        feedback.TicketType(ticketType),
			Title:       title,
			Description: description,
			Reporter:    reporter,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Created #%d: %s\n%s\n", ticket.Number, ticket.Title, ticket.URL)
		return nil
	},
}

func init() {
	submitCmd.Flags().String("app", "", "Source app (rclip, boka, thxbud, mamzo, glasscourt)")
	submitCmd.Flags().String("type", "", "Ticket type (bug, feature-request)")
	submitCmd.Flags().String("title", "", "Ticket title")
	submitCmd.Flags().String("description", "", "Ticket description")
	submitCmd.Flags().String("reporter", "", "Reporter name or email (optional)")
	_ = submitCmd.MarkFlagRequired("app")
	_ = submitCmd.MarkFlagRequired("type")
	_ = submitCmd.MarkFlagRequired("title")
	_ = submitCmd.MarkFlagRequired("description")
}