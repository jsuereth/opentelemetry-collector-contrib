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

An example usage:

```go
ctx := {} // ...create my value...
env := CreateEnvironmentForMyContext()
stmt, err := ParseStatement(env, "ottl(expression)")
// ... error handling ...
result, err := stmt.Execute(context.Background(), ctx)
// ... error handling ...
```

### Creating Custom Context

WIP

### Adding new types

WIP

### Expanding the language

WIP