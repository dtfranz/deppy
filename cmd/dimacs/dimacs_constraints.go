package dimacs

import (
	"context"
	"strings"

	"github.com/operator-framework/deppy/pkg/deppy"
	"github.com/operator-framework/deppy/pkg/deppy/constraint"
	"github.com/operator-framework/deppy/pkg/deppy/input"
)

var _ input.VariableSource = &ConstraintGenerator{}

type ConstraintGenerator struct {
	dimacs *Dimacs
}

func NewDimacsVariableSource(dimacs *Dimacs) *ConstraintGenerator {
	return &ConstraintGenerator{
		dimacs: dimacs,
	}
}

func (d *ConstraintGenerator) GetVariables(_ context.Context) ([]deppy.Variable, error) {
	varMap := make(map[deppy.Identifier]*input.SimpleVariable, len(d.dimacs.variables))
	variables := make([]deppy.Variable, 0, len(d.dimacs.variables))

	for _, id := range d.dimacs.variables {
		variable := input.NewSimpleVariable(deppy.IdentifierFromString(id))
		variables = append(variables, variable)
		varMap[variable.Identifier()] = variable
	}

	// create constraints out of the clauses
	for _, clause := range d.dimacs.clauses {
		terms := strings.Split(clause, " ")
		if len(terms) == 0 {
			continue
		}
		first := terms[0]

		if len(terms) == 1 {
			// TODO: check constraints haven't already been added to the variable
			variable := varMap[deppy.Identifier(strings.TrimPrefix(first, "-"))]
			if strings.HasPrefix(first, "-") {
				variable.AddConstraint(constraint.Not())
			} else {
				// TODO: is this the right constraint here? (given that its an achoring constraint?)
				variable.AddConstraint(constraint.Mandatory())
			}
			continue
		}
		for i := 1; i < len(terms); i++ {
			variable := varMap[deppy.Identifier(strings.TrimPrefix(first, "-"))]
			second := terms[i]
			negSubject := strings.HasPrefix(first, "-")
			negOperand := strings.HasPrefix(second, "-")

			// TODO: this Or constraint is hacky as hell
			variable.AddConstraint(constraint.Or(deppy.Identifier(strings.TrimPrefix(second, "-")), negSubject, negOperand))
			first = second
		}
	}

	return variables, nil
}
