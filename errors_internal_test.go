package k6pack

import (
	"errors"
	"strings"
	"testing"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/require"
)

func Test_checkError(t *testing.T) {
	t.Parallel()

	result := new(api.BuildResult)

	hasError, err := checkError(result)

	require.False(t, hasError)
	require.Nil(t, err)

	result.Errors = append(result.Errors, api.Message{Text: "foo"})

	hasError, err = checkError(result)

	require.True(t, hasError)
	require.NotNil(t, err)
}

func Test_wrapError(t *testing.T) {
	t.Parallel()

	err := wrapError(errors.ErrUnsupported)

	var perr *packError

	require.True(t, errors.As(err, &perr))

	require.Equal(t, "unsupported operation", perr.Error())
}

func Test_packError_Error(t *testing.T) {
	t.Parallel()

	err := packError{messages: []api.Message{
		{
			Text: "Hello, World!",
			Location: &api.Location{
				File:   "foo.ts",
				Column: 42,
				Line:   24,
			},
		},
	}}

	require.Equal(t, "foo.ts:24:42 Hello, World!", err.Error())
}

func Test_packError_Format(t *testing.T) {
	t.Parallel()

	err := wrapError(errors.ErrUnsupported)

	var perr *packError

	require.True(t, errors.As(err, &perr))

	require.True(t, strings.HasSuffix(perr.Format(80, false), "[ERROR] unsupported operation\n\n"))
}
