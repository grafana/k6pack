package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/szkiba/k6pack"
)

func Test_runCmd(t *testing.T) {
	t.Parallel()

	out := filepath.Join(t.TempDir(), "out.js")
	in := filepath.Join("testdata", "simple.ts")

	cmd := newCmd([]string{"-o", out, in})

	require.Equal(t, appname, cmd.Name())
	require.Equal(t, version, cmd.Version)

	runCmd(cmd)

	info, err := os.Stat(out) //nolint:forbidigo
	require.NoError(t, err)

	require.Greater(t, info.Size(), int64(100))
}

func Test_formatError(t *testing.T) {
	t.Parallel()

	str := formatError(errors.ErrUnsupported)

	require.True(t, strings.Contains(str, "[ERROR] unsupported operation\n\n"))

	_, err := k6pack.Imports(`import "foo"`, &k6pack.Options{})
	require.Error(t, err)

	str = formatError(err)

	require.True(t, strings.Contains(str, "[ERROR] Could not resolve"))
}
