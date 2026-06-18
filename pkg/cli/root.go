package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rcli",
	Short: "Unified feedback & ticketing CLI",
	Long:  "Manage feedback tickets across r:clip, BoKa, ThxBud, MamZo, GlassCourt and more.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(feedbackCmd)
}