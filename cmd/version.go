package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// getVersion gets version from environment variable set in main.go
func getVersion() string {
	if v := os.Getenv("SKOOMA_VERSION"); v != "" {
		return v
	}
	return "dev"
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Skooma version information",
	Long:  `Show the Skooma version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Skooma version %s\n", getVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
