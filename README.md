# Package structfields

> Enumerate struct types and their fields from Go source files.

[![Go Reference](https://pkg.go.dev/badge/github.com/icholy/structfields.svg)](https://pkg.go.dev/github.com/icholy/structfields)

**Notes**

* The intended use-case is generating docs for config structs.
* Embedded struct fields are flattened into regular fields.
* Embedded fields are not de-duplicated.
* Fields with inline struct types are not well supported.
