package cli

import "github.com/spf13/cobra"

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Manage feedback tickets",
}

func init() {
	feedbackCmd.AddCommand(submitCmd)
	feedbackCmd.AddCommand(listCmd)
	feedbackCmd.AddCommand(viewCmd)
	feedbackCmd.AddCommand(commentCmd)
	feedbackCmd.AddCommand(updateCmd)
	feedbackCmd.AddCommand(agentsCheckCmd)
}