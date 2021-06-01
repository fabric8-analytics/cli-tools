package maven_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/maven"
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

func (tc isSupportedManifestTestcase) TestMatcher_IsSupportedManifest(t *testing.T) {
	got := tc.Matcher.IsSupportedManifestFormat(tc.Name)
	want := tc.Want
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func (dt depsTreeTestCase) TestDepsTree_GeneratorDependencyTree(t *testing.T) {
	want := dt.Want
	got := dt.Matcher.GeneratorDependencyTree(dt.ManifestFile)
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

// TestMatcher tests the maven matcher.
func TestMatcher(t *testing.T) {
	tt := []isSupportedManifestTestcase{
		{
			Name:    "testdata/requirements.txt",
			Want:    false,
			Matcher: &maven.Matcher{},
		},
		{
			Name:    "testdata/pom.txt",
			Want:    false,
			Matcher: &maven.Matcher{},
		},
		{
			Name:    "testdata/pom.xml",
			Want:    true,
			Matcher: &maven.Matcher{},
		},
		{
			Name:    "testdata/abc.xml",
			Want:    false,
			Matcher: &maven.Matcher{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, tc.TestMatcher_IsSupportedManifest)
	}
}

// TestDepsTree tests the maven Tree Generator.
func TestDepsTree(t *testing.T) {
	tt := []depsTreeTestCase{
		{
			ManifestFile: filepath.Join("testdata", "pom.xml"),
			Want:         filepath.Join(os.TempDir(), "dependencies.txt"),
			Matcher:      &maven.Matcher{},
		},
	}
	for _, dt := range tt {
		t.Run(dt.ManifestFile, dt.TestDepsTree_GeneratorDependencyTree)
	}
}
