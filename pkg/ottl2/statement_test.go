// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"github.com/stretchr/testify/assert"
)

// Define the context we'll evaluate statements against.
type testContext struct {
	name string
}

// Define the structure of the context for the parser to understand.
var testContextType = types.NewStructureType(
	"MyContext",
	map[string]types.Type{
		"name": stdlib.StringType,
	},
)

// Note: MyContext MUST implement types.Val and traits.StructureAccessible
func (m *testContext) ConvertTo(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("unable to convert context")
}

// We define how to access values here.
func (m *testContext) GetField(field string) types.Val {
	switch field {
	case "name":
		return stdlib.NewStringVal(m.name)
	}
	return stdlib.NewErrorVal(fmt.Errorf("unknown field: %s", field))
}
func (m *testContext) Type() types.Type {
	return testContextType
}
func (m *testContext) Value() any {
	return m
}

func IsEmptyFunc() types.Function {
	return stdlib.NewSimpleFunc("IsEmpty", 1, func(args []types.Val) types.Val {
		r, err := args[0].ConvertTo(reflect.TypeFor[string]())
		if err != nil {
			return stdlib.NewErrorVal(err)
		}
		return stdlib.NewBoolVal(len(r.(string)) == 0)
	})
}

func RouteFunc() types.Function {
	return stdlib.NewSimpleFunc("route", 0, func(args []types.Val) types.Val {
		return stdlib.NewBoolVal(true)
	})
}

func Test_simple_e2e(t *testing.T) {
	env := NewTransformContext[testContext](
		testContextType,
		func(v testContext) types.Val { return &v },
		WithFunctions[testContext]([]types.Function{
			IsEmptyFunc(),
			RouteFunc(),
		}),
	)
	stmt, err := ParseStatement(env, "route() where IsEmpty(name)")
	assert.Nil(t, err)
	result, cond, err := stmt.Execute(context.Background(), testContext{name: "test"})
	assert.Nil(t, err)
	assert.False(t, cond)
	assert.Nil(t, result)
}
