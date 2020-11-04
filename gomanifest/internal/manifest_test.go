package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDepsPackages = map[string]GoPackage{
	"github.com/google/go-cmp/cmp": GoPackage{
		Root:       "/home/dhpatel/Documents/code/go-learn/src/go-cmp",
		ImportPath: "github.com/google/go-cmp/cmp",
		Module: GoModule{
			Path:    "github.com/google/go-cmp",
			Main:    true,
			Version: "",
			Replace: nil},
		Standard: false,
		Imports: []string{
			"bytes",
			"fmt",
			"github.com/google/go-cmp/cmp/internal/diff",
			"github.com/google/go-cmp/cmp/internal/flags",
			"github.com/google/go-cmp/cmp/internal/function",
			"github.com/google/go-cmp/cmp/internal/value",
			"math/rand",
			"unicode unicode/utf8",
			"unsafe",
		},
		Deps: []string{
			"bytes",
			"errors",
			"fmt",
			"github.com/google/go-cmp/cmp/internal/diff",
			"github.com/google/go-cmp/cmp/internal/flags",
			"github.com/google/go-cmp/cmp/internal/function",
			"github.com/google/go-cmp/cmp/internal/value",
			"internal/bytealg",
			"internal/cpu",
			"internal/testlog",
			"internal/unsafeheader",
			"io",
			"math/rand",
			"os",
			"reflect",
			"regexp",
			"regexp/syntax",
			"runtime",
			"runtime/internal/sys",
			"unicode/utf8",
			"unsafe",
		},
	},
	"github.com/google/go-cmp/cmp/cmpopts": GoPackage{
		Root:       "/home/dhpatel/Documents/code/go-learn/src/go-cmp",
		ImportPath: "github.com/google/go-cmp/cmp/cmpopts",
		Module: GoModule{
			Path:    "github.com/google/go-cmp",
			Main:    true,
			Version: "",
			Replace: nil,
		},
		Standard: false,
		Imports: []string{
			"fmt",
			"github.com/google/go-cmp/cmp",
			"github.com/google/go-cmp/cmp/internal/function",
			"golang.org/x/xerrors",
			"math",
		},
		Deps: []string{
			"fmt",
			"github.com/google/go-cmp/cmp",
			"github.com/google/go-cmp/cmp/internal/diff",
			"github.com/google/go-cmp/cmp/internal/flags",
			"github.com/google/go-cmp/cmp/internal/function",
			"github.com/google/go-cmp/cmp/internal/value",
			"golang.org/x/xerrors",
			"golang.org/x/xerrors/internal",
			"internal/bytealg",
			"internal/cpu",
		},
	},
	"github.com/google/go-cmp/cmp/internal/diff": GoPackage{
		Root:       "/home/dhpatel/Documents/code/go-learn/src/go-cmp",
		ImportPath: "github.com/google/go-cmp/cmp/internal/diff",
		Module: GoModule{
			Path:    "github.com/google/go-cmp",
			Main:    true,
			Version: "",
			Replace: nil,
		},
		Standard: false,
		Imports: []string{
			"github.com/google/go-cmp/cmp/internal/flags",
			"math/rand",
			"time",
		},
		Deps: []string{
			"errors",
			"github.com/google/go-cmp/cmp/internal/flags",
			"internal/bytealg",
		},
	},
	"github.com/google/go-cmp/cmp/internal/flags": GoPackage{
		Root:       "/home/dhpatel/Documents/code/go-learn/src/go-cmp",
		ImportPath: "github.com/google/go-cmp/cmp/internal/flags",
		Module: GoModule{
			Path:    "github.com/google/go-cmp",
			Main:    true,
			Version: "",
			Replace: nil,
		},
		Standard: false,
		Imports:  []string{},
		Deps:     []string{},
	},
	"golang.org/x/xerrors": GoPackage{
		Root:       "/home/dhpatel/go/pkg/mod/golang.org/x/xerrors@v0.0.0-20191204190536-9bdfabe68543",
		ImportPath: "golang.org/x/xerrors",
		Module: GoModule{
			Path:    "golang.org/x/xerrors",
			Main:    false,
			Version: "v0.0.0-20191204190536-9bdfabe68543",
			Replace: nil,
		},
		Standard: false,
		Imports: []string{
			"fmt",
			"golang.org/x/xerrors/internal",
		},
		Deps: []string{
			"fmt",
			"golang.org/x/xerrors/internal",
			"internal/bytealg",
		},
	},
}

func TestTransformationVerionSemVer(t *testing.T) {
	assert.Equal(t, "2.5.8", transformVersion("2.5.8"), "Semver positive transformation failed")
	assert.Equal(t, "3.2.0", transformVersion("v3.2.0"), "Semver 'v' transformation failed")
	assert.Equal(t, "3.2.0", transformVersion("v3.2.0+incompatible"), "Semver with incompatible transformation failed")
	assert.Equal(t, "3.2.0-alpha", transformVersion("v3.2.0-alpha+incompatible"), "Semver with alpha + incompatible transformation failed")
	assert.Equal(t, "3.2.0-beta1.5", transformVersion("v3.2.0-beta1.5"), "Semver with beta version transformation failed")
	assert.Equal(t, "3.2.0-beta1.2", transformVersion("v3.2.0-beta1.2+incompatible"), "Semver with beta version + incompatible transformation failed")
	assert.Equal(t, "3.2.0-20201023112233-abcd1234abcd", transformVersion("v3.2.0-20201023112233-abcd1234abcd"), "Pseudo version transformation failed")
	assert.Equal(t, "3.2.0-20201023112233-abcd1234abcd", transformVersion("v3.2.0-20201023112233-abcd1234abcd+alpha"), "Pseudo version with alpha transformation failed")
}

func TestBuildManifest(t *testing.T) {
	manifest := BuildManifest(testDepsPackages)
	assert.Equal(t, 1, len(manifest.Packages), "Expected number of deps not found")
}

func TestSaveManifest(t *testing.T) {
	manifestFilePath := testDataFolder + testOutputManifest
	SaveManifestFile(BuildManifest(testDepsPackages), manifestFilePath)

	// Read output json and check for its size
	output := readFileContentForTesting(testOutputManifest)
	assert.Equal(t, 144, len(output), "Output manifest file size missmatch")

	defer os.Remove(manifestFilePath)
}