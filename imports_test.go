package k6pack_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/szkiba/k6pack"
)

func Test_Imports(t *testing.T) {
	t.Parallel()

	imports, err := k6pack.Imports( /*ts*/ `
import "k6"
import "k6/x/foo?bar"
import "k6/x/foo#dummy"
`,
		&k6pack.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"k6", "k6/x/foo#dummy", "k6/x/foo?bar"}, imports)
}

func Test_Imports_error(t *testing.T) {
	t.Parallel()

	_, err := k6pack.Imports( /*ts*/ `
import "foo"
`,
		&k6pack.Options{})

	require.Error(t, err)
}
