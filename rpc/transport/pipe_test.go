package transport_test

import (
	"io"
	"testing"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/rpc/transport"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	t.Parallel()

	m, _ := shark.NewSingleSegmentMessage(nil)

	p1, p2 := transport.NewPipe(1)

	err := p1.Encode(m)
	require.NoError(t, err)

	m2, err := p2.Decode()
	require.NoError(t, err)
	require.NotEqual(t, m, m2, "message should have been copied")

	err = p1.Close()
	require.NoError(t, err)

	err = p1.Encode(m)
	require.ErrorIs(t, err, io.ErrClosedPipe)

	m, err = p2.Decode()
	require.Nil(t, m)
	require.ErrorIs(t, err, io.ErrClosedPipe)
}
