// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import "fmt" // a type to represent paths.
// Values are oneof name, key or index
type testPath struct {
	name  *string
	key   *string
	index *int64
}

func (t testPath) String() string {
	switch {
	// We don't have the ability to depend on traits here, so we simulate it.
	case t.name != nil:
		return fmt.Sprintf(".%s", *t.name)
	case t.key != nil:
		return fmt.Sprintf("[%s]", *t.key)
	case t.index != nil:
		return fmt.Sprintf("[%d]", t.index)
	}
	return "{unknown}"
}

// Helper to lookup a test path.
func lookupTestPath(v Val, path []testPath) Val {
	result := v
	for _, p := range path {
		switch {
		// We don't have the ability to depend on traits here, so we simulate it.
		case p.name != nil:
			result = result.(interface{ GetField(string) Val }).GetField(*p.name)
		case p.key != nil:
			result = result.(interface{ GetKey(string) Val }).GetKey(*p.key)
		case p.index != nil:
			result = result.(interface{ GetIndex(int64) Val }).GetIndex(*p.index)
		}
	}
	return result
}

func fieldPath(f string) testPath {
	return testPath{
		name: &f,
	}
}
func keyPath(k string) testPath {
	return testPath{
		key: &k,
	}
}
func indexPath(i int64) testPath {
	return testPath{
		index: &i,
	}
}

func pathString(path []testPath) string {
	result := "{target}"
	for _, p := range path {
		result += p.String()
	}
	return result
}
