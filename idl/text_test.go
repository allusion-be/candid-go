package idl_test

import "github.com/aviate-labs/candid-go/idl"

func ExampleText() {
	test([]idl.Type{new(idl.TextType)}, []any{""})
	test([]idl.Type{new(idl.TextType)}, []any{"Motoko"})
	test([]idl.Type{new(idl.TextType)}, []any{"Hi ☃\n"})
	// Output:
	// 4449444c00017100
	// 4449444c000171064d6f746f6b6f
	// 4449444c00017107486920e298830a
}
