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

package throttler

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/wesql/wescale/go/vt/topo/topoproto"
)

type rateChange string

const (
	increasedRate rateChange = "increased"
	decreasedRate rateChange = "decreased"
	unchangedRate rateChange = "not changed"
)

type goodOrBadRate string

const (
	goodRate    = "good"
	badRate     = "bad"
	ignoredRate = "ignored"
)

var resultStringTemplate = template.Must(template.New("result.String()").Parse(
	`rate was: {{.RateChange}} from: {{.OldRate}} to: {{.NewRate}}
alias: {{.Alias}} lag: {{.LagRecordNow.Stats.ReplicationLagSeconds}}s
last change: {{.TimeSinceLastRateChange}} rate: {{.CurrentRate}} good/bad? {{.GoodOrBad}} skipped b/c: {{.MemorySkipReason}} good/bad: {{.HighestGood}}/{{.LowestBad}}
state (old/tested/new): {{.OldState}}/{{.TestedState}}/{{.NewState}} 
lag before: {{.LagBefore}} ({{.AgeOfBeforeLag}} ago) rates (primary/replica): {{.PrimaryRate}}/{{.GuessedReplicationRate}} backlog (old/new): {{.GuessedReplicationBacklogOld}}/{{.GuessedReplicationBacklogNew}}
reason: {{.Reason}}`))

// result is generated by the MaxReplicationLag module for each processed
// "replicationLagRecord".
// It captures the details and the decision of the processing.
type result struct {
	Now            time.Time
	RateChange     rateChange
	lastRateChange time.Time

	OldState    state
	TestedState state
	NewState    state

	OldRate int64
	NewRate int64
	Reason  string

	CurrentRate      int64
	GoodOrBad        goodOrBadRate
	MemorySkipReason string
	HighestGood      int64
	LowestBad        int64

	LagRecordNow                 replicationLagRecord
	LagRecordBefore              replicationLagRecord
	PrimaryRate                  int64
	GuessedReplicationRate       int64
	GuessedReplicationBacklogOld int
	GuessedReplicationBacklogNew int
}

func (r result) String() string {
	var b bytes.Buffer
	if err := resultStringTemplate.Execute(&b, r); err != nil {
		panic(fmt.Sprintf("failed to Execute() template: %v", err))
	}
	return b.String()
}

func (r result) Alias() string {
	return topoproto.TabletAliasString(r.LagRecordNow.Tablet.Alias)
}

func (r result) TimeSinceLastRateChange() string {
	if r.lastRateChange.IsZero() {
		return "n/a"
	}
	return fmt.Sprintf("%.1fs", r.Now.Sub(r.lastRateChange).Seconds())
}

func (r result) LagBefore() string {
	if r.LagRecordBefore.isZero() {
		return "n/a"
	}
	return fmt.Sprintf("%ds", r.LagRecordBefore.Stats.ReplicationLagSeconds)
}

func (r result) AgeOfBeforeLag() string {
	if r.LagRecordBefore.isZero() {
		return "n/a"
	}
	return fmt.Sprintf("%.1fs", r.LagRecordNow.time.Sub(r.LagRecordBefore.time).Seconds())
}

// resultRing implements a ring buffer for "result" instances.
type resultRing struct {
	// mu guards the fields below.
	mu sync.Mutex
	// position holds the index of the *next* result in the ring.
	position int
	// wrapped becomes true when the ring buffer "wrapped" at least once and we
	// started reusing entries.
	wrapped bool
	// values is the underlying ring buffer.
	values []result
}

// newResultRing creates a new resultRing.
func newResultRing(capacity int) *resultRing {
	return &resultRing{
		values: make([]result, capacity),
	}
}

// add inserts a new result into the ring buffer.
func (rr *resultRing) add(r result) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.values[rr.position] = r
	rr.position++
	if rr.position == len(rr.values) {
		rr.position = 0
		rr.wrapped = true
	}
}

// latestValues returns all values of the buffer. Entries are sorted in reverse
// chronological order i.e. newer items come first.
func (rr *resultRing) latestValues() []result {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	start := rr.position - 1
	if start == -1 {
		// Current position is at the end.
		start = len(rr.values) - 1
	}
	count := len(rr.values)
	if !rr.wrapped {
		count = rr.position
	}

	results := make([]result, count)
	for i := 0; i < count; i++ {
		pos := start - i
		if pos < 0 {
			// We started in the middle of the array and need to wrap around at the
			// beginning of it.
			pos += count
		}
		results[i] = rr.values[pos%count]
	}
	return results
}
