package npm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/npm"
)

type isSupportedManifestTestcase struct {
	Name    string
	Want    bool
	Matcher driver.StackAnalysisInterface
}

type depsTreeTestCase struct {
	ManifestFile string
	Want         string
	Matcher      driver.StackAnalysisInterface
}

func (tc isSupportedManifestTestcase) TestIsSupportedManifest(t *testing.T) {
	got := tc.Matcher.IsSupportedManifestFormat(tc.Name)
	want := tc.Want
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func (dt depsTreeTestCase) TestGeneratorDependencyTree(t *testing.T) {
	want := dt.Want
	got := dt.Matcher.GeneratorDependencyTree(dt.ManifestFile)
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

// TestMatcher tests the npm matcher.
func TestMatcher(t *testing.T) {
	tt := []isSupportedManifestTestcase{
		{
			Name:    "testdata/requirements.txt",
			Want:    false,
			Matcher: &npm.Matcher{},
		},
		{
			Name:    "testdata/pom.txt",
			Want:    false,
			Matcher: &npm.Matcher{},
		},
		{
			Name:    "testdata/package.json",
			Want:    true,
			Matcher: &npm.Matcher{},
		},
		{
			Name:    "testdata/abc.xml",
			Want:    false,
			Matcher: &npm.Matcher{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, tc.TestIsSupportedManifest)
	}
}

// TestDepsTree tests the npm Tree Generator.
func TestDepsTree(t *testing.T) {
	testdata, _ := filepath.Abs("testdata")
	tt := []depsTreeTestCase{
		{
			ManifestFile: filepath.Join(testdata, "package.json"),
			Want:         filepath.Join(os.TempDir(), "npmlist.json"),
			Matcher:      &npm.Matcher{},
		},
	}
	for _, dt := range tt {
		t.Run(dt.ManifestFile, dt.TestGeneratorDependencyTree)
	}
}
