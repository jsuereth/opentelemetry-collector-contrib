// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

type errorType struct{}

func (e errorType) Name() string {
	return "error"
}

var ErrorType types.Type = errorType{}

type errorVal struct {
	e error
}

// Type implements Val.
func (e errorVal) Type() types.Type {
	return ErrorType
}

func (e errorVal) ConvertTo(t types.Type) (any, error) {
	return nil, e.e
}

func (e errorVal) Value() any {
	return e.e
}

func (e errorVal) String() string {
	return e.e.Error()
}

// We can return errors within values.
func NewErrorVal(e error) types.Val {
	return errorVal{e}
}
