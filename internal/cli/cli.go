// Package cli contains CLI command for k6pack.
package cli

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/szkiba/k6pack"
)

const defaultTimeout = 30 * time.Second

// New creates new cobra command for k6pack command.
func New() *cobra.Command {
	opts := new(k6pack.Options)

	var output string

	cmd := &cobra.Command{
		Use:   "pack [flags] filename",
		Short: "TypeScript transpiler and module bundler for k6.",
		Long:  "Bundle TypeScript/JavaScript sources into a single k6 test script.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			return pack(args[0], opts, out)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()

	cwd, _ := os.Getwd() //nolint:forbidigo

	flags.BoolVar(&opts.TypeScript, "typescript", false, "force TypeScript loader")
	flags.BoolVar(&opts.SourceMap, "sourcemap", false, "emit the source map with an inline data URL")
	flags.StringVar(&opts.SourceRoot, "source-root", cwd, "sets the sourceRoot field in generated source maps")
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

	script, err := k6pack.Pack(string(contents), opts)
	if err != nil {
		return err
	}

	_, err = out.Write(script)

	return err
}
