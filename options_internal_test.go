package k6pack

import (
	"path/filepath"
	"testing"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/require"
)

func Test_options_setDefaults(t *testing.T) {
	t.Parallel()

	opts := new(Options)

	opts.setDefaults()

	require.False(t, opts.SourceMap)
	require.False(t, opts.Minify)
	require.False(t, opts.TypeScript)
	require.Nil(t, opts.Externals)
	require.Equal(t, ".", opts.Directory)
	require.Empty(t, opts.Filename)
	require.Zero(t, opts.Timeout)
	require.Equal(t, api.LoaderJS, opts.loaderType())
	require.Equal(t, api.SourceMapNone, opts.sourceMapType())

	opts.Filename = filepath.Join("foo", "bar.ts")
	opts.Directory = ""

	opts.setDefaults()
	require.Equal(t, "foo", opts.Directory)
	require.True(t, opts.TypeScript)
	require.Equal(t, api.LoaderTS, opts.loaderType())

	opts.SourceMap = true
	require.Equal(t, api.SourceMapInline, opts.sourceMapType())
}

func Test_options_stdinOptions(t *testing.T) {
	t.Parallel()

	opts := new(Options)
	opts.Filename = filepath.Join("foo", "bar.ts")
	opts.setDefaults()

	stdinopts := opts.stdinOptions("Hello, World!")

	require.Equal(t, &api.StdinOptions{
		Contents:   "Hello, World!",
		ResolveDir: "foo",
		Sourcefile: filepath.Join("foo", "bar.ts"),
		Loader:     api.LoaderTS,
	}, stdinopts)
}
