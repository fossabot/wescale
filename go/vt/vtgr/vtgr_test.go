package vtgr

import (
	"context"
	"syscall"
	"testing"
	"time"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/topo/memorytopo"
	"vitess.io/vitess/go/vt/vtgr/config"
	"vitess.io/vitess/go/vt/vtgr/controller"
	"vitess.io/vitess/go/vt/vtgr/db"
	"vitess.io/vitess/go/vt/vttablet/tmclient"

	"github.com/stretchr/testify/assert"
)

func TestSighupHandle(t *testing.T) {
	ctx := context.Background()
	ts := memorytopo.NewServer("cell1")
	defer ts.Close()
	ts.CreateKeyspace(ctx, "ks", &topodatapb.Keyspace{})
	ts.CreateShard(ctx, "ks", "0")
	vtgr := newVTGR(
		ctx,
		ts,
		tmclient.NewTabletManagerClient(),
	)
	var shards []*controller.GRShard
	config := &config.VTGRConfig{
		DisableReadOnlyProtection:   false,
		GroupSize:                   5,
		MinNumReplica:               3,
		BackoffErrorWaitTimeSeconds: 10,
		BootstrapWaitTimeSeconds:    10 * 60,
	}
	shards = append(shards, controller.NewGRShard("ks", "0", nil, vtgr.tmc, vtgr.topo, db.NewVTGRSqlAgent(), config, *localDbPort))
	vtgr.Shards = shards
	shard := vtgr.Shards[0]
	shard.LockShard(ctx, "test")
	res := 0
	vtgr.handleSignal(func(i int) {
		res = i
	})
	assert.NotNil(t, shard.GetUnlock())
	assert.False(t, vtgr.stopped.Get())
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, res)
	assert.Nil(t, shard.GetUnlock())
	assert.True(t, vtgr.stopped.Get())
}
