package groups

import (
	"context"
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
)

const (
	rolesRego = "roles"
	envsRego  = "envs"
	teamsRego = "teams"
	groupRego = "group"
	regoQueryName = "data.test.rbac.allow"
)

func toInput(roles, envs, teams []string) map[string]interface{} {
	return map[string]interface{}{
		rolesRego: roles,
		envsRego:  envs,
		teamsRego: teams,
	}
}

func prepareQueries(rules map[string]string) map[string]*rego.PreparedEvalQuery {
	preparedQueries := make(map[string]*rego.PreparedEvalQuery, len(rules))
	for name, rule := range rules {
		preparedQueries[name] = preparedQuery(name, rule)
	}

	return preparedQueries
}

func preparedQuery(name, rule string) *rego.PreparedEvalQuery {
	ctx := context.Background()
	query, err := rego.New(
		rego.Query(regoQueryName),
		rego.Module(name+".rego", rule),
	).PrepareForEval(ctx)
	if err != nil {
		log.Print(fmt.Errorf("prepare query for '%s': %w", name, err))
		return nil
	}

	return &query
}
