package main

import (
	"github.com/harukitosa/merucari/errchecker"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(errchecker.Analyzer) }

