package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rcli/feedback/internal/feedback"
)

var agentsCheckCmd = &cobra.Command{
	Use:   "agents-check",
	Short: "Summarize open tickets for AI agent context",
	Example: `  rcli feedback agents-check --app rclip`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, _ := cmd.Flags().GetString("app")
		if app == "" {
			return fmt.Errorf("--app is required")
		}

		client := NewAPIClient()
		tickets, err := client.List(feedback.ListFilter{
			App:   app,
			State: "open",
		})
		if err != nil {
			return err
		}

		var b strings.Builder
		fmt.Fprintf(&b, "# Open feedback for %s (%d tickets)\n\n", app, len(tickets))

		for _, t := range tickets {
			fmt.Fprintf(&b, "## #%d — %s [%s/%s]\n", t.Number, t.Title, t.Type, t.Status)
			fmt.Fprintf(&b, "%s\n", t.Body)
			fmt.Fprintf(&b, "URL: %s\n\n", t.URL)
		}

		fmt.Fprint(os.Stdout, b.String())
		return nil
	},
}

func init() {
	agentsCheckCmd.Flags().String("app", "", "App to check (rclip, boka, thxbud, mamzo, glasscourt)")
	_ = agentsCheckCmd.MarkFlagRequired("app")
}