package candid

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

// Definition represents an imports or type definition.
type Definition interface {
	def()
	fmt.Stringer
}

func (t Type) def()   {}
func (i Import) def() {}

// Type represents a named type definition.
type Type struct {
	Id   string
	Data Data
}

func (t Type) String() string {
	return fmt.Sprintf("type %s = %s", t.Id, t.Data.String())
}

func convertType(n *ast.Node) Type {
	var (
		id   = n.FirstChild
		data = n.LastChild
	)
	return Type{
		Id:   id.Value,
		Data: convertData(data),
	}
}

// Import represents an import declarations from another file.
type Import struct {
	Text string
}

func (i Import) String() string {
	return fmt.Sprintf("import %q", i.Text)
}
