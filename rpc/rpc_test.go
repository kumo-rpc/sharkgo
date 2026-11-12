package rpc

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/rpc/internal/testshark"
)

func TestConnection_BaseContext(t *testing.T) {
	t.Parallel()

	t.Run("background context", func(t *testing.T) {
		client, server := net.Pipe()
		doneCh := make(chan struct{}, 1)

		go func() {
			bootstrapClient := testshark.StreamTest_ServerToClient(slowStreamTestServer{})
			conn := NewConn(NewStreamTransport(server), &Options{
				BootstrapClient: shark.Client(bootstrapClient),
			})
			defer conn.Close()

			<-conn.Done()
			close(doneCh)
		}()

		func() {
			conn := NewConn(NewStreamTransport(client), nil)
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := testshark.StreamTest(conn.Bootstrap(ctx))
			defer client.Release()

			err := client.Push(ctx, func(st testshark.StreamTest_push_Params) error {
				return st.SetData(make([]byte, 1))
			})

			require.NoError(t, err)
		}()

		<-doneCh
	})

	t.Run("external context", func(t *testing.T) {
		client, server := net.Pipe()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		clientConnectedCh := make(chan struct{}, 1)
		contextCancelledCh := make(chan struct{}, 1)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			bootstrapClient := testshark.StreamTest_ServerToClient(slowStreamTestServer{})
			conn := NewConn(NewStreamTransport(server), &Options{
				BootstrapClient: shark.Client(bootstrapClient),
				BaseContext: func() context.Context {
					return ctx
				},
			})
			defer conn.Close()

			select {
			case <-clientConnectedCh:
				cancel()
				contextCancelledCh <- struct{}{}
			case <-conn.Done():
				t.Failed()
			}

			// Connection should close when external context is cancelled
			<-conn.Done()
		}()

		go func() {
			defer wg.Done()

			conn := NewConn(NewStreamTransport(client), nil)
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := testshark.StreamTest(conn.Bootstrap(ctx))
			defer client.Release()

			if err := client.Resolve(ctx); err != nil {
				require.NoError(t, err)
			}

			clientConnectedCh <- struct{}{}
			<-contextCancelledCh

			err := client.Push(ctx, func(st testshark.StreamTest_push_Params) error {
				return st.SetData(make([]byte, 1))
			})

			require.NoError(t, err)
		}()

		wg.Wait()
	})
}
