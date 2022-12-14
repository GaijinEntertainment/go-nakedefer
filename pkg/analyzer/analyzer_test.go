package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/GaijinEntertainment/go-nakedefer/pkg/analyzer"
)

func TestAll(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	a, err := analyzer.NewAnalyzer(
		[]string{"ignoreFunc", "os\\.(Create|WriteFile|Chmod)", "fmt\\.Print.*", "io\\.Close"},
	)
	if err != nil {
		t.Error(err)
	}

	analysistest.Run(t, testdata, a, "p")
}
