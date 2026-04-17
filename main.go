// CLI to scaffold fullstack SPAs with Go, TypeScript, React, Tailwind, and Vite.
package main

import (
	"os"

	"github.com/skooma-cli/skooma/cmd"
	"github.com/skooma-cli/skooma/internal/config"
)

var version = "0.2.0"

func main() {
	os.Setenv("SKOOMA_VERSION", version)

	// Load config to ensure it exists and is valid before executing any commands
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
