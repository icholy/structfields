package testmod

// E is a struct
//go:what
//go:directive
type E struct {
	// F1 is a string
	F1 string `json:"Test"` // F1 is the first field

	F2 int // F2 only has a comment

	// F3 only has a doc
	F3 bool
}
