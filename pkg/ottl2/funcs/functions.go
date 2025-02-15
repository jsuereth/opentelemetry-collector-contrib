// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

// Standard functions for OTTL.
func StandardFuncs() []runtime.Function {
	return []runtime.Function{
		NewIsBoolFunc(),
		NewIsIntFunc(),
		NewIsDoubleFunc(),
		NewIsStringFunc(),
		NewIsListFunc(),
		NewIsMapFunc(),
		NewSetFunc(),
		NewAppendFunc(),
	}
}
