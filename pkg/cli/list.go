package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/rcli/feedback/internal/feedback"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List feedback tickets",
	Example: `  rcli feedback list --app rclip --status open
  rcli feedback list --app all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, _ := cmd.Flags().GetString("app")
		status, _ := cmd.Flags().GetString("status")
		ticketType, _ := cmd.Flags().GetString("type")
		state, _ := cmd.Flags().GetString("state")

		if app == "" {
			app = "all"
		}

		client := NewAPIClient()
		tickets, err := client.List(feedback.ListFilter{
			App:    app,
			Status: status,
			Type:   ticketType,
			State:  state,
		})
		if err != nil {
			return err
		}

		if len(tickets) == 0 {
			fmt.Println("No tickets found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "#\tAPP\tTYPE\tSTATUS\tTITLE")
		for _, t := range tickets {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				t.Number, t.App, t.Type, t.Status, t.Title)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().String("app", "all", "Filter by app (or 'all')")
	listCmd.Flags().String("status", "open", "Filter by status (open, triaged, in-progress, done)")
	listCmd.Flags().String("type", "", "Filter by type (bug, feature-request)")
	listCmd.Flags().String("state", "open", "GitHub state (open, closed, all)")
}