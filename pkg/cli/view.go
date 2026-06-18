package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [number]",
	Short: "View a feedback ticket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var number int
		if _, err := fmt.Sscanf(args[0], "%d", &number); err != nil {
			return fmt.Errorf("invalid issue number: %s", args[0])
		}

		client := NewAPIClient()
		ticket, err := client.Get(number)
		if err != nil {
			return err
		}

		fmt.Printf("#%d %s\n", ticket.Number, ticket.Title)
		fmt.Printf("App:     %s\n", ticket.App)
		fmt.Printf("Type:    %s\n", ticket.Type)
		fmt.Printf("Status:  %s\n", ticket.Status)
		fmt.Printf("State:   %s\n", ticket.State)
		fmt.Printf("URL:     %s\n", ticket.URL)
		fmt.Printf("Labels:  %v\n", ticket.Labels)
		fmt.Printf("\n%s\n", ticket.Body)
		return nil
	},
}