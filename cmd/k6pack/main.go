// Package main contains the k6pack CLI entry point.
package main

import (
	"log"
	"os"

	"github.com/grafana/k6pack/cmd"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	appname = "k6pack"
	version = "dev"
)

func main() {
	runCmd(newCmd(os.Args[1:])) //nolint:forbidigo
}

func newCmd(args []string) *cobra.Command {
	cmd := cmd.New()
	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	log.SetFlags(0)
	log.Writer()

	if err := cmd.Execute(); err != nil {
		log.Fatal(formatError(err))
	}
}
