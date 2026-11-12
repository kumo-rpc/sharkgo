package rpc_test

import (
	"context"
	"testing"

	"github.com/kumo-rpc/sharkgo"
	"github.com/kumo-rpc/sharkgo/rpc/internal/testshark"
	"github.com/stretchr/testify/assert"
)

func TestPipelineChaining(t *testing.T) {
	t.Parallel()

	bg := context.Background()

	np := func(cat testshark.CapArgsTest_self_Params) error {
		return nil
	}

	ctx, cancel := context.WithCancel(bg)
	cli := testshark.CapArgsTest_ServerToClient(&delayingCapReturner{ctx})
	defer cli.Release()

	// Originally, chaining pipelines like this didn't work correctly.

	fut1, rel1 := cli.Self(bg, np)
	defer rel1()
	// This one was fine
	fut2, rel2 := fut1.Self().Self(bg, np)
	defer rel2()
	// This one would be delivered to the wrong place
	// (the same place as the previous one)
	fut3, rel3 := fut2.Self().Self(bg, np)
	defer rel3()
	// This one would segfault
	fut4, rel4 := fut3.Self().Self(bg, np)
	defer rel4()

	cancel()
	_, err := fut4.Struct()

	assert.Nil(t, err)
}

type delayingCapReturner struct {
	context.Context
}

func (c delayingCapReturner) Call(ctx context.Context, call testshark.CapArgsTest_call) error {
	return shark.Unimplemented("yes")
}

func (c delayingCapReturner) Self(ctx context.Context, call testshark.CapArgsTest_self) error {
	call.Go()
	select {
	case <-ctx.Done():
	case <-c.Done():
	}

	results, err := call.AllocResults()
	if err != nil {
		return err
	}

	results.SetSelf(testshark.CapArgsTest_ServerToClient(selfCapReturner{}))
	return nil
}

type selfCapReturner struct{}

func (c selfCapReturner) Call(ctx context.Context, call testshark.CapArgsTest_call) error {
	return shark.Unimplemented("yes")
}

func (c selfCapReturner) Self(ctx context.Context, call testshark.CapArgsTest_self) error {
	results, err := call.AllocResults()
	if err != nil {
		return err
	}

	results.SetSelf(testshark.CapArgsTest_ServerToClient(c))
	return nil
}
