// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package otlpcel

import (
	"github.com/google/cel-go/cel"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type Statement struct {
	filter cel.Program
	setter cel.Program
}

func CompileOtlpCel(filter string, setter string) (*Statement, error) {
	env, err := NewSpanEnv()
	if err != nil {
		return nil, err
	}
	fast, iss := env.Compile(filter)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	fp, err := env.Program(fast)
	if err != nil {
		return nil, err
	}
	sast, iss := env.Compile(setter)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	sp, err := env.Program(sast)
	if err != nil {
		return nil, err
	}
	return &Statement{fp, sp}, nil
}

func (s *Statement) Eval(span ptrace.Span) (bool, any, error) {
	activation, err := NewSpanActivation(span)
	if err != nil {
		return false, nil, err
	}
	out, _, err := s.filter.Eval(activation)
	if err != nil {
		return false, nil, err
	}
	cond := out.Value().(bool)
	if cond {
		result, _, err := s.setter.Eval(activation)
		if err != nil {
			return true, nil, err
		}
		return true, result.Value(), nil
	}
	return false, nil, nil
}
