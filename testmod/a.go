package testmod

import (
	"time"

	"testmod/other" // different
	"testmod/sub"
)

type A struct {
	F1     string
	F2     int
	F3     bool
	F4     time.Duration
	F5, F6 byte
}

type B struct {
	sub.C
	different.X
}
