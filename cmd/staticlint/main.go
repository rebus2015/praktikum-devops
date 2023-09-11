package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/rebus2015/praktikum-devops/cmd/staticlint/analyzers"
)

func main() {
	multichecker.Main(
		analyzers.Get()...,
	)
}
