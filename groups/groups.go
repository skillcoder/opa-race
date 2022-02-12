package groups

import (
	"context"

	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

var (
	errRuleNotFound                  = errors.New("rule not found")
	errEmptyRules                    = errors.New("empty rules")
	errUnexpectedRegoExpressionValue = errors.New("unexpected rego expressions value")
	errDenied                        = errors.New("denied")
)

type Groups struct {
	preparedQueries map[string]*rego.PreparedEvalQuery
}

func New(rules map[string]string) *Groups {
	preparedQueries := prepareQueries(rules)

	return &Groups{
		preparedQueries: preparedQueries,
	}
}

func (g *Groups) Check(ctx context.Context, name string, roles, envs, teams []string) error {
	query, err := g.getQuery(name)
	if err != nil {
		return errors.Wrapf(err, "getting query for '%s'", name)
	}

	res, err := query.Eval(ctx, rego.EvalInput(toInput(roles, envs, teams)))
	if err != nil {
		return errors.Wrapf(err, "eval rego code for '%s'", name)
	}

	if len(res) == 0 {
		return errEmptyRules
	}

	value, ok := res[0].Expressions[0].Value.([]interface{})
	if !ok {
		return errors.Wrapf(errUnexpectedRegoExpressionValue, "type '%T'", res[0].Expressions[0].Value)
	}

	if len(value) == 0 {
		return errDenied
	}

	return nil
}

func (g *Groups) getQuery(name string) (*rego.PreparedEvalQuery, error) {
	if query, ok := g.preparedQueries[name]; ok {
		return query, nil
	}

	return nil, errRuleNotFound
}
