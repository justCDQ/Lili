package examplesgo

import (
	"slices"
	"sync/atomic"
)

type Snapshot struct {
	RouteNames []string
}

type Registry struct {
	requests atomic.Uint64
	snapshot atomic.Pointer[Snapshot]
}

func (r *Registry) RecordRequest() {
	r.requests.Add(1)
}

func (r *Registry) RequestCount() uint64 {
	return r.requests.Load()
}

func (r *Registry) Publish(routes []string) {
	copyOfRoutes := slices.Clone(routes)
	r.snapshot.Store(&Snapshot{RouteNames: copyOfRoutes})
}

func (r *Registry) Routes() []string {
	snapshot := r.snapshot.Load()
	if snapshot == nil {
		return nil
	}
	return slices.Clone(snapshot.RouteNames)
}
