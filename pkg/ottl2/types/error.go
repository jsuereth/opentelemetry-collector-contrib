// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import "reflect" // Allows errors returned from evaluating expressions

type errorType struct{}

func (e errorType) Name() string {
	return "error"
}

var ErrorType Type = errorType{}

type errorVal struct {
	e error
}

// Type implements Val.
func (e errorVal) Type() Type {
	return ErrorType
}

func (e errorVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	return nil, e.e
}

func (e errorVal) Value() any {
	return e.e
}

// We can return errors within values.
func NewErrorVal(e error) Val {
	return errorVal{e}
}
