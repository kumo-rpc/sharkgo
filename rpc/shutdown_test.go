package rpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/rpc"
	"github.com/kumo-rpc/sharkgo/rpc/internal/testshark"
	"github.com/kumo-rpc/sharkgo/rpc/transport"
)

// TestRejectOnDisconnect verifies that, when a connection is dropped, outstanding calls
// fail with a "disconnected" exception.
func TestRejectOnDisconnect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	serverNetConn, clientNetConn := net.Pipe()

	// The remote server will close this to acknowledge that the call has reached the other
	// side; this makes sure we are testing the right logic, not the scenario where the
	// connection is dropped before the message is sent, which could fail with an IO error,
	// as opposed to being rejected on shutdown.
	readyCh := make(chan struct{})

	serverRpcConn := rpc.NewConn(transport.NewStream(serverNetConn), &rpc.Options{
		BootstrapClient: shark.Client(testshark.PingPong_ServerToClient(dropPingServer{
			readyCh: readyCh,
		})),
	})

	clientRpcConn := rpc.NewConn(transport.NewStream(clientNetConn), nil)

	client := testshark.PingPong(clientRpcConn.Bootstrap(ctx))
	defer client.Release()
	future, release := client.EchoNum(ctx, nil)
	defer release()
	<-readyCh
	serverNetConn.Close()
	_, err := future.Struct()
	if !shark.IsDisconnected(err) {
		t.Fatalf("Wanted disconnected error but got %v", err)
	}
	<-serverRpcConn.Done()
	<-clientRpcConn.Done()
}

type dropPingServer struct {
	readyCh chan<- struct{}
}

func (s dropPingServer) EchoNum(ctx context.Context, p testshark.PingPong_echoNum) error {
	p.Go()
	close(s.readyCh)
	<-ctx.Done()
	return nil
}
