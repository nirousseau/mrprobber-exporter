package main

import (
	"mrprober/cmd"
)

// ----------------------------------------------------------------------------

// VERSION This var is overridden at build time with make
var VERSION = "SNAPSHOT"

// ----------------------------------------------------------------------------

func main() {
	cmd.Execute(VERSION)
}
