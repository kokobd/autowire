# Autowire

Google golang [wire](https://github.com/google/wire) is an excellent compile
-time dependency injection tool.

However, with wire, you have to manually write `wire.go`, which can be
 tedious and error-prune.
 
Autowire build on top of wire, and generates `wire.go` and `wire_gen.go` for
you. All you have to do is to use annotations to mark your component, just
like what you will do in Spring Boot.