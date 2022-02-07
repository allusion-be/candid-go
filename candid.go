package candid

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/aviate-labs/candid-go/did"
	"github.com/aviate-labs/candid-go/idl"
	"github.com/aviate-labs/candid-go/internal/candid"
	"github.com/aviate-labs/candid-go/internal/candidvalue"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
)

func DecodeValue(value []byte) (string, error) {
	types, values, err := idl.Decode(value)
	if err != nil {
		return "", err
	}
	if len(types) != 1 || len(values) != 1 {
		return "", fmt.Errorf("can not decode: %x", value)
	}
	s, err := valueToString(types[0], values[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s)", s), nil
}

func DecodeValues(types []idl.Type, values []interface{}) (string, error) {
	var ss []string
	if len(types) != len(values) {
		return "", fmt.Errorf("unequal length")
	}
	for i := range types {
		s, err := valueToString(types[i], values[i])
		if err != nil {
			return "", err
		}
		ss = append(ss, s)
	}
	return fmt.Sprintf("(%s)", strings.Join(ss, "; ")), nil
}

func EncodeValue(value string) ([]byte, error) {
	p, err := ast.New([]byte(value))
	if err != nil {
		return nil, err
	}
	n, err := candidvalue.Values(p)
	if err != nil {
		return nil, err
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		return nil, err
	}
	types, args, err := did.ConvertValues(n)
	if err != nil {
		return nil, err
	}
	return idl.Encode(types, args)
}

// ParseDID parses the given raw .did files and returns the Program that is defined in it.
func ParseDID(raw []byte) (did.Description, error) {
	p, err := ast.New(raw)
	if err != nil {
		return did.Description{}, err
	}
	n, err := candid.Prog(p)
	if err != nil {
		return did.Description{}, err
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		return did.Description{}, err
	}
	return did.ConvertDescription(n), nil
}

func valueToString(typ idl.Type, value interface{}) (string, error) {
	switch t := typ.(type) {
	case *idl.Opt:
		if value == nil {
			return "opt null", nil
		}
		s, err := valueToString(t.Type, value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("opt %s", s), nil
	case *idl.Nat:
		if t.Base == 0 {
			return fmt.Sprintf("%s : nat", value), nil
		}
		return fmt.Sprintf("%s : nat%d", value, t.Base), nil
	case *idl.Int:
		if t.Base == 0 {
			return fmt.Sprintf("%s", value), nil
		}
		return fmt.Sprintf("%s : int%d", value, t.Base), nil
	case *idl.Float:
		f, _ := value.(*big.Float).Float64()
		return fmt.Sprintf("%.f : float%d", f, t.Base), nil
	case *idl.Bool:
		return fmt.Sprintf("%t", value), nil
	case *idl.Null:
		return "null", nil
	case *idl.Text:
		return fmt.Sprintf("%q", value), nil
	case *idl.Rec:
		var ss []string
		for _, f := range t.Fields {
			v := value.(map[string]interface{})
			s, err := valueToString(f.Type, v[f.Name])
			if err != nil {
				return "", nil
			}
			ss = append(ss, fmt.Sprintf("%s = %s", f.Name, s))
		}
		if len(ss) == 0 {
			return "record {}", nil
		}
		return fmt.Sprintf("record { %s }", strings.Join(ss, "; ")), nil
	case *idl.Variant:
		f := t.Fields[0]
		v := value.(*idl.FieldValue).Value
		var s string
		switch t := f.Type.(type) {
		case *idl.Null:
			s = f.Name
		default:
			sv, err := valueToString(t, v)
			if err != nil {
				return "", err
			}
			s = fmt.Sprintf("%s = %s", f.Name, sv)
		}
		return fmt.Sprintf("variant { %s }", s), nil
	case *idl.Vec:
		var ss []string
		for _, a := range value.([]interface{}) {
			s, err := valueToString(t.Type, a)
			if err != nil {
				return "", err
			}
			ss = append(ss, s)
		}
		if len(ss) == 0 {
			return "vec {}", nil
		}
		return fmt.Sprintf("vec { %s }", strings.Join(ss, "; ")), nil
	default:
		panic(fmt.Sprintf("%s, %v", typ, value))
	}
}
