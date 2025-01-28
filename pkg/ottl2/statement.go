// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"reflect"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// Statement holds a top level Statement for processing telemetry data. A Statement is a combination of a function
// invocation and the boolean expression to match telemetry for invoking the function.
type Statement[K EvalContext] struct {
	function Interpretable
	// TODO - we probably want a better interpreatable that gives us a boolean
	condition         Interpretable
	origText          string
	telemetrySettings component.TelemetrySettings
}

func (s Statement[K]) Execute(ctx context.Context, env K) (any, bool, error) {
	condition, err := s.condition.Eval(ctx, env).ConvertTo(reflect.TypeFor[bool]())
	defer func() {
		if s.telemetrySettings.Logger != nil {
			s.telemetrySettings.Logger.Debug("TransformContext after statement execution", zap.String("statement", s.origText), zap.Bool("condition matched", condition.(bool)), zap.Any("TransformContext", env))
		}
	}()
	if err != nil {
		return nil, false, err
	}
	var result any
	if condition.(bool) {
		result = s.function.Eval(ctx, env).Value()
		// TODO - faster error checking.
		if reflect.TypeOf(result).AssignableTo(reflect.TypeFor[error]()) {
			return nil, true, result.(error)
		}
	}
	return result, condition.(bool), nil
}
