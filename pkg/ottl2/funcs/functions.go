// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

// Standard functions for OTTL.
func StandardFuncs() []types.Function {
	return []types.Function{
		NewIsIntFunc(),
		NewSetFunc(),
	}
}
