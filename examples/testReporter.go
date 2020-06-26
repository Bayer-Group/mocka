package examples

import (
	"log"
	"os"
)

// t is used in place of *t.Testing for examples
var t = log.New(os.Stderr, "example: ", log.Lshortfile)
