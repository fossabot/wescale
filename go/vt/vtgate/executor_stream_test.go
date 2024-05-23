/*
Copyright ApeCloud, Inc.
Licensed under the Apache v2(found in the LICENSE file in the root directory).
*/

/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vtgate

import (
	"testing"
	"time"

	querypb "github.com/wesql/wescale/go/vt/proto/query"

	"github.com/wesql/wescale/go/cache"
	"github.com/wesql/wescale/go/vt/discovery"
	topodatapb "github.com/wesql/wescale/go/vt/proto/topodata"

	"context"

	"github.com/stretchr/testify/require"

	"github.com/wesql/wescale/go/sqltypes"
	_ "github.com/wesql/wescale/go/vt/vtgate/vindexes"
	"github.com/wesql/wescale/go/vt/vttablet/sandboxconn"
)

func TestStreamSQLSharded(t *testing.T) {
	cell := "aa"
	hc := discovery.NewFakeHealthCheck(nil)
	s := createSandbox("TestExecutor")
	s.VSchema = executorVSchema
	getSandbox(KsTestUnsharded).VSchema = unshardedVSchema
	serv := newSandboxForCells([]string{cell})
	resolver := newTestResolver(hc, serv, cell)
	shards := []string{"-20", "20-40", "40-60", "60-80", "80-a0", "a0-c0", "c0-e0", "e0-"}
	for _, shard := range shards {
		_ = hc.AddTestTablet(cell, shard, 1, "TestExecutor", shard, topodatapb.TabletType_PRIMARY, true, 1, nil)
	}
	executor := NewExecutor(context.Background(), serv, cell, resolver, false, false, testBufferSize, cache.DefaultConfig, nil, false, querypb.ExecuteOptions_V3)

	sql := "stream * from sharded_user_msgs"
	result, err := executorStreamMessages(executor, sql)
	require.NoError(t, err)
	wantResult := &sqltypes.Result{
		Fields: sandboxconn.SingleRowResult.Fields,
		Rows: [][]sqltypes.Value{
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
			sandboxconn.StreamRowResult.Rows[0],
		},
	}
	if !result.Equal(wantResult) {
		t.Errorf("result: %+v, want %+v", result, wantResult)
	}
}

func executorStreamMessages(executor *Executor, sql string) (qr *sqltypes.Result, err error) {
	results := make(chan *sqltypes.Result, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err = executor.StreamExecute(
		ctx,
		"TestExecuteStream",
		NewSafeSession(primarySession),
		sql,
		nil,
		func(qr *sqltypes.Result) error {
			results <- qr
			return nil
		},
	)
	close(results)
	if err != nil {
		return nil, err
	}
	first := true
	for r := range results {
		if first {
			qr = &sqltypes.Result{Fields: r.Fields}
			first = false
		}
		qr.Rows = append(qr.Rows, r.Rows...)
	}
	return qr, nil
}
