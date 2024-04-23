package cli

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/grafana/k6pack"
	"github.com/stretchr/testify/require"
)

func Test_pack_error(t *testing.T) {
	t.Parallel()

	opts := new(k6pack.Options)

	err := pack("no_such_file", opts, io.Discard)

	require.Error(t, err)

	err = pack(filepath.Join("testdata", "invalid_script.js"), opts, io.Discard)

	require.Error(t, err)
}
