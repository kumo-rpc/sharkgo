package schemas_test

import (
	"testing"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/internal/schema"
	"github.com/kumo-rpc/sharkgo/schemas"
	gocp "github.com/kumo-rpc/sharkgo/std/go"
)

func TestDefaultFind(t *testing.T) {
	gocp.RegisterSchema(schemas.DefaultRegistry)
	if s := schemas.Find(0xdeadbeef); s != nil {
		t.Errorf("schemas.Find(0xdeadbeef) = %d-byte slice; want nil", len(s))
	}
	s := schemas.Find(gocp.Package_)
	if s == nil {
		t.Fatalf("schemas.Find(%#x) = nil", gocp.Package_)
	}
	msg, err := shark.Unmarshal(s)
	if err != nil {
		t.Fatalf("shark.Unmarshal(schemas.Find(%#x)) error: %v", gocp.Package_, err)
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		t.Fatalf("ReadRootCodeGeneratorRequest error: %v", err)
	}
	nodes, err := req.Nodes()
	if err != nil {
		t.Fatalf("req.Nodes() error: %v", err)
	}
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		if n.Id() == gocp.Package_ {
			// Found
			if n.Which() != schema.Node_Which_annotation {
				t.Errorf("found node %#x which = %v; want annotation", gocp.Package_, n.Which())
			}
			return
		}
	}
	t.Fatalf("could not find node %#x in registry", gocp.Package_)
}

func TestNotFound(t *testing.T) {
	reg := new(schemas.Registry)
	_, err := reg.Find(0)
	if err == nil {
		t.Error("new(schemas.Registry).Find(0) = nil; want not found error")
	}
	if !schemas.IsNotFound(err) {
		t.Errorf("new(schemas.Registry).Find(0) = %v; want not found error", err)
	}
}
