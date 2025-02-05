// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

type EnumProvider interface {
	TypeName() string
	FindName(id int64) (string, bool)
	ResolveName(name string) (int64, bool)
}
