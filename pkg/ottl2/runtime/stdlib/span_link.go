// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime" // TODO - structural type
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanLinkType runtime.Type = runtime.NewPrimitiveType("SpanLink")

func NewSpanLinkVar(link ptrace.SpanLink) runtime.Var {
	panic("unimplemented")
}
