package groups_test

import (
	"context"
	"sync"
	"testing"

	"opa-race/groups"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var rules = `package test.rbac

allow[{"group": "group1", "hash": "test-hash"}] {
	input.roles[_] == "role1"
	lower(input.teams[_]) == "team1"
}

allow[{"group": "group2", "hash": "test-hash"}] {
	input.roles[_] == "role2"
	lower(input.envs[_]) == ["e1", "e2"][_]
}
`

// TestGroup_CheckParallel detect race conditions in rego.Eval
//  race in opa hash in v0.27.1 and even in v0.37.2
//  https://github.com/open-policy-agent/opa/issues/2129
func TestGroup_CheckParallel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rulesMap := map[string]string{
		"g1": rules,
	}

	ag := groups.New(rulesMap)

	roles := []string{"role1", "role2", "role3"}
	envs := []string{"env1", "env2", "env3"}
	teams := []string{"team1", "team2", "team3"}

	// NOTE: uncomment and race will have despaired
	/*
		err := ag.Check(
			context.TODO(),
			"g1",
			roles,
			envs,
			teams,
		)
		assert.NoError(t, err)
		if err != nil {
			t.FailNow()
		}
	*/

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(t *testing.T) {
			t.Helper()

			defer wg.Done()

			err := ag.Check(
				context.TODO(),
				"g1",
				roles,
				envs,
				teams,
			)
			assert.NoError(t, err)
		}(t)
	}

	wg.Wait()
}
