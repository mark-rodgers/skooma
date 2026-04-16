// CLI to scaffold fullstack SPAs with Go, TypeScript, React, Tailwind, and Vite.
package main

import (
	"os"

	"github.com/skooma-cli/skooma/cmd"
	"github.com/skooma-cli/skooma/internal/config"
)

var version = "0.1.0-dev"

// main is the entry point for the CLI application.
func main() {
	os.Setenv("SKOOMA_VERSION", version)

	// Load config to ensure it exists and is valid before executing any commands
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
