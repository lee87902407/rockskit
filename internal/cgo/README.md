# internal/cgo

This directory is reserved for the future CGO bridge.

All future `import "C"`, `#cgo`, pointer conversion, native lifetime tracking, and `char** errptr` handling should live here.
No Go source files are added yet so the current prototype keeps its existing build shape.
