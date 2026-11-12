package server_test

import (
	"fmt"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/server"
)

func ExampleIsServer() {
	x := int(42)
	c := shark.NewClient(server.New([]server.Method{}, x, nil))
	snapshot := c.Snapshot()
	brand := snapshot.Brand()
	snapshot.Release()
	if brand, ok := server.IsServer(brand); ok {
		fmt.Println("Client is a server, got brand:", brand)
	}
	// Output:
	// Client is a server, got brand: 42
}
