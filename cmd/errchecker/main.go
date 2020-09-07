package main

import (
	"github.com/harukitosa/errchecker"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(errchecker.Analyzer) }
