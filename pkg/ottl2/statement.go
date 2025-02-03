// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// Statement holds a top level Statement for processing telemetry data. A Statement is a combination of a function
// invocation and the boolean expression to match telemetry for invoking the function.
type Statement[K any] struct {
	function Interpretable
	// TODO - we probably want a better interpreatable that gives us a boolean
	condition         Interpretable
	origText          string
	ctx               TransformContext[K]
	telemetrySettings component.TelemetrySettings
}

func (s Statement[K]) String() string {
	return fmt.Sprintf("ottl[%s]", s.origText)
}

// Helper to access logger so we can shift it to different input later.
func (s Statement[K]) logger() *zap.Logger {
	return s.telemetrySettings.Logger
}

func (s Statement[K]) Execute(ctx context.Context, env *K) (any, bool, error) {
	realEnv := s.ctx.NewEvalContext(env)
	condition, err := s.condition.Eval(ctx, realEnv).ConvertTo(stdlib.BoolType)
	defer func() {
		if s.logger() != nil {
			s.logger().Debug("TransformContext after statement execution", zap.String("statement", s.origText), zap.Bool("condition matched", condition.(bool)), zap.Any("TransformContext", env))
		}
	}()
	if err != nil {
		return nil, false, err
	}
	var result any
	if condition.(bool) {
		result = s.function.Eval(ctx, realEnv).Value()
		// TODO - faster error checking.
		if result != nil && reflect.TypeOf(result).AssignableTo(reflect.TypeFor[error]()) {
			return nil, true, result.(error)
		}
		// panic(fmt.Sprintf("evaluated %v, using %v", env, realEnv))
	}
	return result, condition.(bool), nil
}

// Parses an OTTL statement for a given context, returning something that can evaluate
// against that context.
func ParseStatement[R any](ctx TransformContext[R], statement string) (Statement[R], error) {
	p := NewParser(ctx)
	parsed, err := parseRawStatement(statement)
	if err != nil {
		return Statement[R]{}, err
	}
	condition, expr, err := p.parseStatement(*parsed)
	if err != nil {
		return Statement[R]{}, err
	}
	return Statement[R]{
		function:          expr,
		condition:         condition,
		origText:          statement,
		ctx:               ctx,
		telemetrySettings: component.TelemetrySettings{}, // TODO - move these to context
	}, nil
}
