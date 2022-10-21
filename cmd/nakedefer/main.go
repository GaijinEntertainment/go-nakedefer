package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/GaijinEntertainment/go-nakedefer/pkg/analyzer"
)

func main() {
	a, err := analyzer.NewAnalyzer([]string{})
	if err != nil {
		panic(err)
	}

	singlechecker.Main(a)
}
