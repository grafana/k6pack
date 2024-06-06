// Package main contains CLI documentation generator tool.
package main

import (
	_ "embed"
	"strings"

	"github.com/grafana/clireadme"
	"github.com/grafana/k6pack/cmd"
)

func main() {
	root := cmd.New()
	root.Use = strings.ReplaceAll(root.Use, "pack", "k6pack")
	clireadme.Main(root, 1)
}
