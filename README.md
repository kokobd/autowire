# Autowire

*WARNING: CURRENTLY IN DEVELOPMENT*

Google golang [wire](https://github.com/google/wire) is an excellent compile
-time dependency injection tool.

However, with wire, you have to manually write `wire.go`, which can be
 tedious and error-prune.
 
Autowire build on top of wire, and generates `wire.go` and `wire_gen.go` for
you. All you have to do is to use annotations to mark your component, just
like what you will do in Spring Boot.

## User Guide

### Start watcher

Run `autowire -watch` from your project root (where `go.mod` is located), then `wire.go` should be generated on the fly

### Add annotations

These examples covers most common use cases.

#### Function provider
Simply annotate your constructor with `@Component`

```golang
// @Component
func NewConfg() Config {
    // ...
}
```
This generates:
```golang
wire.Build(
    // ...
    pkg.NewConfig,
    // ...
)
```

#### Struct provider with all fields injected

```golang
// @Component(allFields = true)
// @Bind(value = IMyAwesomeService, pointer = true)
type MyAwesomeService struct {
    // ...
}
```

This generates
```golang
wire.Build(
    // ...
    wire.Struct(pkg.MyAwesomeService, "*")
    wire.Bind(new(pkg.IMyAwesomeService), new(*pkg.MyAwesomeService))
    // ...
)
```

#### Struct provider with selected fields
We omit @Bind in this example for simplicity, but you can add it if you want.
```golang
// @Component
type MyAwesomeService struct {
    // @Autowire
    AwesomeConfig AwesomeConfig
}
```

This generates:
```golang
wire.Build(
    // ...
    wire.Struct(pkg.MyAwesomeService, "AwesomeConfig")
    // ...
)
```

#### Custom provider set

You can specify a provider set directly, in case you have a component not supported by existing features of `autowire`.

```golang
// @Component
var CustomProviderSet = wire.NewSet(
    // ... your providers
)
```
This generates
```golang
wire.Build(
    // ...
    CustomProviderSet,
    // ...
)
```

## Reference

### Command Line

- `-watch` automatically detect file change, then regenerate wire.go on the fly
- `-out` specify where to write output. default: `cmd/wire.go`

### Annotations

#### Component

Scope: function or struct

Params:
- `allFields` (default `false`) whether to inject all public fields

#### Autowire

Scope: struct field

Params: none

#### Bind

Scope: struct

Params:
- `value` (required) the interface to bind
- `pointer` (default `true`) if set true, pointer to the struct is used as implementation of the specified interface. Otherwise, value of the struct is used.
