package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment [number] [message]",
	Short: "Add a comment to a ticket",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var number int
		if _, err := fmt.Sscanf(args[0], "%d", &number); err != nil {
			return fmt.Errorf("invalid issue number: %s", args[0])
		}

		message := args[1]
		if len(args) > 2 {
			message = joinArgs(args[1:])
		}

		client := NewAPIClient()
		if err := client.Comment(number, message); err != nil {
			return err
		}

		fmt.Printf("Comment added to #%d\n", number)
		return nil
	},
}

func joinArgs(args []string) string {
	result := args[0]
	for _, a := range args[1:] {
		result += " " + a
	}
	return result
}