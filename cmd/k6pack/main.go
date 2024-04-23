// Package main contains the k6pack CLI entry point.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/grafana/k6pack/internal/cli"
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
	cmd := cli.New()
	cmd.Use = strings.Replace(cmd.Use, cmd.Name(), appname, 1)
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
