# Why a rewrite?

OTTL has proven highly useful as a technology for collector processing.  It finds its way into many
collector receivers beyond the transform processor, due to a few key capabilities:

- Simple boolean expression logic
- Simple transformation capabilities ("editors")
- Easily embedded within YAML configuration.

However, the current internal architecture is reaching limitations that prevent expansion 
of the language. OTTL is hoping to expand in a few key ways:

- Supporting collections
  - [Determine an approach to looping](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/29289) *Note: this is the initial push for this proposal*
  - [Check if an element is in a collection](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/30420)
  - [Support dynamic/embedded context, e.g. for group_by capability](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/34649)
  - [Improve index validation for slice vs. map](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/37646)
- Consolidating type coercion rules
  - [Map literals cannot be set in slices](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/37405)
- Static analysis capabilities
  - [Inferring context from OTTL expression](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/37904)
  - [Optimising execution via the OTLP hierarchy](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/29016)
  - [Expanding OTTL into other domains](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/25161)

## Fundamental limitation

The current internal architecture of OTTL limits growth due to a simple architectural decision:

All expressions and statements are tied to a specific context type (see [statement.go](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/parser.go#L20)).

This makes simple nested contexts (e.g. for loops, let expressions, lambda expressions) practically impossible (or at least, incredibly difficult) as you need the ability to push/pop a context *stack* for evaluating.

For example, when encountering a pythonic expression like `[x * 2 for x in context.list]`, the interpretre would create a new context for every element in the list and evaluate the expression `x*2`, where `x` looks up the current element in the list.

In OTTL today, this is very hard and involves a myriad of casting and "hope" that you haven't broken the getter/setter convention against a specific context.  While OTTL *does* allow creating a new `Context` which has a "higher" context it delegates too, the Parse needs the ability to do this on the fly, with for-expressions, lambda expressions, etc.  Additionally, the interpreter needs to push/pop these contexts, requiring access to do so in [`Expr`](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/expression.go#L24).

This problem is not insurmountable, but they continue to carry a high amount of friction for language expansion.

## Larger limitations

While the notion of nested context is, in my opinion, enough to warrant a reworking of the core of OTTL, there are
further limitations that limit OTTL's growth.

Briefly, these are:

- Terms (names you reference from context) being hidden within getter/setter machinery, preventing static analysis.
- Functions having access to context can prevent understanding whether an expression only uses Resource or requires access to Span  (`IsRootSpan()` is highly problematic here).   This prevents fixing [#29016](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/29016)

## What can/should we improve?

We need to give OTTL the right foundation so we can grow it over time to meet the needs of the ecosystem.  This includes:

  - Terms (context and paths available on context, e.g. `span.status`) should be known to the Parser.
  - Types can be ascribed to terms on context, and validated (e.g. preventing integer index to a map)
  - Expressions should be decomposible such that understanding if a condition needs access to a signal, scope or resource is something that can be statically determined.
  - 

## Design

See the [design doc](design.md) for this rework.  However, to re-iterate key components:

- Three tiered approach to OTTL
  - Syntax Tree: existing grammar.go file. Raw syntax nodes parsed from OTTL statements.
  - Expression Tree: an *interactable* view of the syntax tree, with types resolved such that you can perform analyiss and optimisations.
  -  Value/Runtime layer: An abstraction for evaluating OTTL expressions.  This lets us standardize type conversions, list&map expressions, etc.

Within this new design:

- Creating a new type for OTTL is done via Value/Runtime layer.
- Creating a new `Function` for OTTL would interact with the Value/Runtime layer.
- Creating a new `Context` for OTTL would interact with the Value/Runtime layer, and can continue re-use exisitng contexts.  This would look similar to adding any new type to OTTL.
- Optimising expressions would interact with the Expression layer.  It should be possible to understand which components of context are required to evaluate an expression.

Since the value/runtime layer is most important, here's a simple example of what it would include:

- runtime type information  (Type, StructureType,Indexable)
- runtime value abstractions (Var, Val)

```go
type Type interface {
    Name() string
}
type StructureType interface {
    Type
    GetField(name string) (Type, bool)
}
type Val interface {
    Type() Type
    Value() any
    ConvertTo(Type) (any, error)
}
type Var interface {
    Val
    Set(Val) error
}
type Indexable interface {
    Val
    GetIndex(int64) Val
    Len() int64
}
type KeyIndexable interface {
    Val
    GetKey(string) Val
    Len() int64
}
type Structure interface {
    Val
    GetField(string) Val
}
```

When parsing an expression, you simply need a `map[string]Type` available to understand what values the context provides.

When evaluating an expression, you simply need a `map[string]Val` available to manipulate values.

OTTL would provide standard `Val` wrappers for all types needed in OTLP.
