package marshal

import "github.com/aviate-labs/candid-go/idl"

type Context[T idl.Type] struct {
	tdt *idl.TypeDefinitionTable
	typ T
}

func NewContext() Context[idl.Type] {
	return Context[idl.Type]{
		tdt: &idl.TypeDefinitionTable{
			Indexes: make(map[string]int),
		},
	}
}

func NewContextWithType[T idl.Type](t T) Context[T] {
	return Context[T]{
		tdt: &idl.TypeDefinitionTable{
			Indexes: make(map[string]int),
		},
		typ: t,
	}
}

func ContextToType[T idl.Type, M idl.Type](ctx Context[T], t M) Context[M] {
	return Context[M]{
		tdt: ctx.tdt,
		typ: t,
	}
}
