package candidtest_test

import (
	"io/ioutil"
	"testing"

	"github.com/aviate-labs/candid-go/internal/candidtest"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
)

func TestData(t *testing.T) {
	rawDid, _ := ioutil.ReadFile("../../idl/testdata/prim.test.did")
	p, _ := ast.New(rawDid)
	if _, err := candidtest.TestData(p); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		t.Error(err)
	}
}
