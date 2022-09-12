package did

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/candid-go/internal/candidvalue"
	"github.com/di-wu/parser/ast"
)

func ConvertValues(n *ast.Node) ([]idl.Type, []interface{}, error) {
	switch n.Type {
	case candidvalue.BoolValueT:
		switch n.Value {
		case "true":
			return []idl.Type{new(idl.BoolType)}, []interface{}{true}, nil
		case "false":
			return []idl.Type{new(idl.BoolType)}, []interface{}{false}, nil
		default:
			panic(n)
		}
	case candidvalue.NullT:
		return []idl.Type{new(idl.NullType)}, []interface{}{nil}, nil
	case candidvalue.NumT:
		typ, arg, err := convertNum(n)
		if err != nil {
			return nil, nil, err
		}
		return []idl.Type{typ}, []interface{}{arg}, nil
	case candidvalue.OptValueT:
		types, args, err := ConvertValues(n.Children()[0])
		if err != nil {
			return nil, nil, err
		}
		return []idl.Type{idl.NewOptionalType(types[0])}, []interface{}{args[0]}, nil
	case candidvalue.RecordT:
		if len(n.Children()) == 0 {
			return []idl.Type{idl.NewRecordType(nil)}, []interface{}{nil}, nil
		}
		types := make(map[string]idl.Type)
		args := make(map[string]interface{})
		for _, n := range n.Children() {
			n := n.Children()
			id := n[0].Value
			typ, arg, err := ConvertValues(n[1])
			if err != nil {
				return nil, nil, err
			}
			types[id] = typ[0]
			args[id] = arg[0]
		}
		return []idl.Type{idl.NewRecordType(types)}, []interface{}{args}, nil
	case candidvalue.TextT:
		n := n.Children()[0]
		s := strings.TrimPrefix(strings.TrimSuffix(n.Value, "\""), "\"")
		return []idl.Type{new(idl.TextType)}, []interface{}{s}, nil
	case candidvalue.ValuesT:
		var (
			types []idl.Type
			args  []interface{}
		)
		for _, n := range n.Children() {
			idl, arg, err := ConvertValues(n)
			if err != nil {
				return nil, nil, err
			}
			types = append(types, idl...)
			args = append(args, arg...)
		}
		return types, args, nil
	case candidvalue.VariantT:
		n := n.Children()
		id := n[0].Value
		switch len(n) {
		case 1:
			typ := idl.NewVariantType(map[string]idl.Type{id: new(idl.NullType)})
			arg := idl.Variant{Name: id, Value: nil, Type: typ}
			return []idl.Type{typ}, []interface{}{arg}, nil
		case 2:
			varType, varArg, err := ConvertValues(n[1])
			if err != nil {
				return nil, nil, err
			}
			typ := idl.NewVariantType(map[string]idl.Type{id: varType[0]})
			arg := idl.Variant{Name: id, Value: varArg[0], Type: typ}
			return []idl.Type{typ}, []interface{}{arg}, nil
		default:
			panic(n)
		}
	case candidvalue.VecT:
		n := n.Children()
		if len(n) == 0 {
			return []idl.Type{idl.NewVectorType(new(idl.NullType))}, []interface{}{[]interface{}{}}, nil
		}
		var types idl.Type
		var args []interface{}
		for _, n := range n {
			typ, arg, err := ConvertValues(n)
			if err != nil {
				return nil, nil, err
			}
			types = typ[0]
			args = append(args, arg[0])
		}
		return []idl.Type{idl.NewVectorType(types)}, []interface{}{args}, nil
	default:
		panic(n)
	}
}

func convertNum(n *ast.Node) (idl.Type, any, error) {
	switch n := n.Children(); len(n) {
	case 1:
		n := n[0]

		// float64
		if strings.Contains(n.Value, ".") {
			v := strings.ReplaceAll(n.Value, "_", "")
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, nil, err
			}
			return idl.Float64Type(), f, nil
		}

		// int
		v := strings.ReplaceAll(n.Value, "_", "")
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		return new(idl.IntType), big.NewInt(i), nil
	case 2:
		vArg := n[0].Value
		vType := n[1].Value

		// floats
		if vType == "float32" || vType == "float64" {
			v := strings.ReplaceAll(vArg, "_", "")
			switch n := n[1]; n.Value {
			case "float32":
				f, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return nil, nil, err
				}
				return idl.Float32Type(), float32(f), nil
			default:
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, nil, err
				}
				return idl.Float64Type(), f, nil
			}
		}

		// ints
		v := strings.ReplaceAll(vArg, "_", "")
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		switch vType {
		case "nat":
			bi := big.NewInt(i)
			return new(idl.NatType), bi, nil
		case "nat8":
			return idl.Nat8Type(), uint8(i), nil
		case "nat16":
			return idl.Nat16Type(), uint16(i), nil
		case "nat32":
			return idl.Nat32Type(), uint32(i), nil
		case "nat64":
			return idl.Nat64Type(), uint64(i), nil
		case "int":
			bi := big.NewInt(i)
			return new(idl.IntType), bi, nil
		case "int8":
			return idl.Int8Type(), int8(i), nil
		case "int16":
			return idl.Int16Type(), int16(i), nil
		case "int32":
			return idl.Int32Type(), int32(i), nil
		case "int64":
			return idl.Int64Type(), i, nil
		default:
			panic(n)
		}
	default:
		panic(n)
	}
}
