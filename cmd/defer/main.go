package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/GaijinEntertainment/go-defer/pkg/analyzer"
)

func main() {
	a, err := analyzer.NewAnalyzer()
	if err != nil {
		panic(err)
	}

	singlechecker.Main(a)
}
