package testutil

import (
	"time"
	"testing"
	"github.com/hashicorp/consul/consul/structs"
)

type testFn func() (bool, error)
type errorFn func(error)

func WaitForResult(test testFn, error errorFn) {
	retries := 100  // 5 seconds timeout

	for retries > 0 {
		time.Sleep(50 * time.Millisecond)
		retries--

		success, err := test()
		if success {
			return
		}

		if retries == 0 {
			error(err)
		}
	}
}

type clientRPC func(string, interface {}, interface {}) error

func WaitForLeader(t *testing.T, rpc clientRPC, args interface{}) {
	WaitForResult(func() (bool, error) {
		var out structs.IndexedNodes
		err := rpc("Catalog.ListNodes", args, &out)
		return out.QueryMeta.KnownLeader, err
	}, func(err error) {
		t.Fatalf("failed to find leader: %v", err)
	})
}
