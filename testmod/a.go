package testmod

import (
	"time"

	"testmod/sub"
)

type A struct {
	F1 string
	F2 int
	F3 bool
	F4 time.Duration
}

type B struct {
	A
	sub.C
}
