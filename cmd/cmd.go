// Package cmd contains cobra command for k6pack.
package cmd

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/grafana/k6pack"
	"github.com/spf13/cobra"
)

const defaultTimeout = 30 * time.Second

//go:embed help.md
var help string

// New creates new cobra command for "k6pack" command.
func New() *cobra.Command {
	opts := new(k6pack.Options)

	var output string

	cmd := &cobra.Command{
		Use:   "k6pack [flags] filename",
		Short: "TypeScript transpiler and module bundler for k6.",
		Long:  help,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var out io.Writer

			if len(output) == 0 {
				out = os.Stdout //nolint:forbidigo
			} else {
				file, err := os.Create(filepath.Clean(output)) //nolint:forbidigo
				if err != nil {
					return err
				}

				defer file.Close() //nolint:errcheck

				out = file
			}

			if opts.SourceRoot == "." {
				abs, err := filepath.Abs(opts.SourceRoot)
				if err != nil {
					return err
				}

				opts.SourceRoot = abs
			}

			return pack(args[0], opts, out)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()

	flags.BoolVar(&opts.TypeScript, "typescript", false, "force TypeScript loader")
	flags.BoolVar(&opts.SourceMap, "sourcemap", false, "emit the source map with an inline data URL")
	flags.StringVar(&opts.SourceRoot, "source-root", ".", "sets the sourceRoot field in generated source maps")
	flags.BoolVar(&opts.Minify, "minify", false, "minify the output")
	flags.DurationVar(&opts.Timeout, "timeout", defaultTimeout, "HTTP timeout for remote modules")
	flags.StringVarP(&output, "output", "o", "", "write output to file (default stdout)")
	flags.StringArrayVar(&opts.Externals, "external", opts.Externals, "exclude module(s) from the bundle")

	return cmd
}

func pack(filename string, opts *k6pack.Options, out io.Writer) error {
	opts.Filename = filename

	fname, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(fname) //nolint:forbidigo,gosec
	if err != nil {
		return err
	}

	script, _, err := k6pack.Pack(string(contents), opts)
	if err != nil {
		return err
	}

	_, err = out.Write(script)

	return err
}
