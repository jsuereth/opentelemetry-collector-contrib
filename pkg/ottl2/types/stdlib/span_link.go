// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types" // TODO - structural type
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanLinkType types.Type = types.NewPrimitiveType("SpanLink")

func NewSpanLinkVar(link ptrace.SpanLink) types.Var {
	panic("unimplemented")
}
