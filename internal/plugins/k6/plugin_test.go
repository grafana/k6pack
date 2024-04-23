package k6_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/grafana/k6pack/internal/plugins/k6"
	"github.com/stretchr/testify/require"
)

func Test_plugin_add_k6_to_external(t *testing.T) {
	t.Parallel()

	script := /*ts*/ `
import { sleep } from "k6"
			
export default function() {
	sleep(1)
}
`
	result := pack(t, script, false)
	require.Empty(t, result.Errors)
}

func Test_plugin_metafile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		imports  []string
		wantErr  bool
	}{
		{name: "single_k6", location: "k6", imports: []string{"k6"}},
		{name: "extension", location: "k6/x/foo", imports: []string{"k6/x/foo"}},
		{name: "extension_with_prefix", location: "xk6-foo", imports: []string{"k6/x/foo", "xk6-foo"}},
		{name: "extension_with_query", location: "k6/x/foo?answer=42", imports: []string{"k6/x/foo?answer=42"}},
		{name: "extension_with_hash", location: "k6/x/foo#bar", imports: []string{"k6/x/foo#bar"}},
		{name: "extension_with_scope", location: "@john/xk6-foo", imports: []string{"@john/xk6-foo", "k6/x/foo"}},
		{name: "extension_with_scope_and_query", location: "@john/xk6-foo?answer=42", imports: []string{"@john/xk6-foo?answer=42", "k6/x/foo?answer=42"}},
		{name: "extension_with_scope_and_path_and_query", location: "@john/xk6-foo/sub1?answer=42", imports: []string{"@john/xk6-foo/sub1?answer=42", "k6/x/foo/sub1?answer=42"}},
		{name: "extension_with_scope_and_paths_and_query", location: "@john/xk6-foo/sub1/sub2?answer=42", imports: []string{"@john/xk6-foo/sub1/sub2?answer=42", "k6/x/foo/sub1/sub2?answer=42"}},
		// error
		{name: "invalid_import", location: ":%4!?", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			script := /*ts*/ `
			import "` + tt.location + `"
			
			export default function() {
			}
			`

			result := pack(t, script, true)

			if tt.wantErr {
				require.NotEmpty(t, result.Errors)
				return
			}

			require.Empty(t, result.Errors)

			var meta metaFile

			err := json.Unmarshal([]byte(result.Metafile), &meta)
			require.NoError(t, err)

			require.NotNil(t, meta.K6)

			require.EqualValues(t, tt.imports, meta.K6.Imports)

			fmt.Println(meta.K6)
		})
	}
}

type metaFile struct {
	K6 *k6.Metadata `json:"k6,omitempty"`
}

func pack(t *testing.T, script string, metafile bool) api.BuildResult {
	t.Helper()

	const httpBase = "http://127.0.0.1"

	if strings.Contains(script, httpBase) {
		server := httptest.NewServer(http.FileServer(http.Dir("testdata")))
		defer server.Close()

		script = strings.ReplaceAll(script, httpBase, server.URL)
	}

	return api.Build(api.BuildOptions{
		Bundle: true,
		Stdin: &api.StdinOptions{
			Contents:   script,
			ResolveDir: ".",
			Sourcefile: t.Name(),
			Loader:     api.LoaderTS,
		},
		LogLevel: api.LogLevelSilent,
		Plugins:  []api.Plugin{k6.New()},
		Metafile: metafile,
	})
}
