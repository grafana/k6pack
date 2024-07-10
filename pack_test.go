package k6pack_test

import (
	"fmt"
	"testing"

	"github.com/grafana/k6pack"
	"github.com/stretchr/testify/require"
)

func Test_Pack_error(t *testing.T) {
	t.Parallel()

	_, _, err := k6pack.Pack( /*ts*/ `
import "foo"
`,
		&k6pack.Options{})

	require.Error(t, err)
}

func Test_Pack(t *testing.T) {
	t.Parallel()

	src, meta, err := k6pack.Pack( /*ts*/ `
import { User, newUser } from "./examples/user"
const user : User = newUser("John")
console.log(user)
`, &k6pack.Options{TypeScript: true})

	fmt.Println(string(src))

	require.NoError(t, err)

	exp := /*js*/ `var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);

// examples/user.ts
var UserAccount = class {
  constructor(name) {
    __publicField(this, "name");
    __publicField(this, "id");
    this.name = name;
    this.id = Math.floor(Math.random() * Number.MAX_SAFE_INTEGER);
  }
};
function newUser(name) {
  return new UserAccount(name);
}

// <stdin>
var user = newUser("John");
console.log(user);
`

	require.Equal(t, exp, string(src))
	require.NotNil(t, meta)
	require.Empty(t, meta.Imports)
}

func Test_Pack_Imports(t *testing.T) {
	t.Parallel()

	_, meta, err := k6pack.Pack( /*ts*/ `
import "k6"
import "k6/x/foo?bar"
import "k6/x/foo#dummy"
`,
		&k6pack.Options{})

	require.NoError(t, err)
	require.Equal(t, []string{"k6", "k6/x/foo#dummy", "k6/x/foo?bar"}, meta.Imports)
}

func Test_Pack_Imports_error(t *testing.T) {
	t.Parallel()

	_, _, err := k6pack.Pack( /*ts*/ `
import "foo"
`,
		&k6pack.Options{})

	require.Error(t, err)
}
