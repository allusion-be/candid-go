package idl_test

import (
	"math/big"

	"github.com/aviate-labs/candid-go/idl"
)

func ExampleOpt() {
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{nil})
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{big.NewInt(1)})
	// Output:
	// 4449444c016e7d010000
	// 4449444c016e7d01000101
}
