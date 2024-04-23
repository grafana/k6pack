package cli_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/k6pack/internal/cli"
	"github.com/stretchr/testify/require"
)

func Test_New_stdout(t *testing.T) {
	t.Parallel()

	cmd := cli.New()

	args := []string{filepath.Join("testdata", "simple.ts")}

	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)

	err := cmd.Execute()

	require.NoError(t, err)
}

func Test_New_file(t *testing.T) {
	t.Parallel()

	cmd := cli.New()

	filename := filepath.Join(t.TempDir(), "out.js")

	args := []string{"-o", filename, filepath.Join("testdata", "simple.ts")}

	cmd.SetArgs(args)

	err := cmd.Execute()
	require.NoError(t, err)

	info, err := os.Stat(filename) //nolint:forbidigo
	require.NoError(t, err)

	require.Greater(t, info.Size(), int64(100))
}

func Test_New_file_error(t *testing.T) {
	t.Parallel()

	cmd := cli.New()

	args := []string{"-o", filepath.Join("no_such_dir", "no_such_subdir", "file.js"), filepath.Join("testdata", "simple.ts")}

	cmd.SetArgs(args)

	err := cmd.Execute()
	require.Error(t, err)
}
