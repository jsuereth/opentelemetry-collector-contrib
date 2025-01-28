# OTTL Design

This document outlines the architecture of the OTTL internals reboot.

## Architecture

The OTTL package is split into several layers, some of which are internal only:

- **Syntax Tree Layer**: 
- **Expression Layer**: A layer that provides nodes which can be interpreted using a context.
  - Includes the `Interpretable` interface and `Statement[R]` struct.
  - See `expr.go` for details.
- **Value Layer**: A layer over runtime types that provides casting and 
  reflection components necessary for execution.
  - Includes the `Val` and `Var` interfaces, for defining context.
  - See [types](types/README.md) package

### Entry Points

All of this is wrapped in a few "starter" points for using OTTL.

- `ParseStatement()`: This parses an OTTL statement and returns a `Statement[R]`.
  A `Statement[R]` can be evalauted, repeatedly, against a context `R`.
- `TransformContext[E any]` is a structure you need to construct for any valid OTTL context.

Here's an example using a custom context:

```go
env := NewTransformContext[MyContext](
  MyContextType, 
  func(v testContext) types.Val { return &v },
  WithFunction[MyContext]("IsEmpty", IsEmptyFunc()),
  WithFunction[MyContext]("route", RouteFunc()),
)
stmt, err := ParseStatement(env, "route() where IsEmpty(name)")
// ... error handling ...
result, cond, err := stmt.Execute(context.Background(), MyContext { name: "test"})
// result: nil
// cond: false
// err: nil
```

### Creating Custom Context

Creating a new context:

```go
// Define the context we'll evaluate statements against.
type MyContext struct { 
   name string
}

// Define the structure of the context for the parser to understand.
var MyContextType = types.NewStructureType(
  "MyContext",
  map[string]types.Type {
    "name": types.StringType
  }
)

// Note: MyContext MUST implement types.Val and traits.StructureAccessible
func (m *MyContext) ConvertTo(typeDesc reflect.Type) (any, error) { ... }
// We define how to access values here.
func (m *MyContext) GetField(field string) types.Val {
  switch field {
  case "name":
    return types.NewStringVal(m.name)
  }
  return types.NewErrorVal(fmt.Errorf("unknown field: %s", field))
}
func (m *MyContext) Type() types.Type { 
  return MyContextType 
}
func (m *MyContext) Value() any { 
  return m
}
```

### Adding new functions

WIP

### Adding new types

WIP

### Expanding the language

WIP