package idl_test

import (
	"github.com/aviate-labs/candid-go/idl"
)

func ExampleNull() {
	test([]idl.Type{new(idl.NullType)}, []any{nil})
	// Output:
	// 4449444c00017f
}
