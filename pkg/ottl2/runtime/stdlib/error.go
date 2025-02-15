// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

type errorType struct{}

func (e errorType) Name() string {
	return "error"
}

var ErrorType runtime.Type = errorType{}

type errorVal struct {
	e error
}

// Type implements Val.
func (e errorVal) Type() runtime.Type {
	return ErrorType
}

func (e errorVal) ConvertTo(t runtime.Type) (any, error) {
	return nil, e.e
}

func (e errorVal) Value() any {
	return e.e
}

func (e errorVal) String() string {
	return e.e.Error()
}

// We can return errors within values.
func NewErrorVal(e error) runtime.Val {
	return errorVal{e}
}
