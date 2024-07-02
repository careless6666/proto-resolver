package functional_tests

import (
	"ProtoDepsResolver/cmd/app"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestUrlOnlyDeps(t *testing.T) {
	// Arrange

	depsFile := `version: v1
deps:
  - url: https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto github.com/googleapis/googleapis/google/api v1`

	testDir, err := createdepsFileAndDir(t, depsFile)

	err = os.Chdir(testDir)
	if err != nil {
		require.NoError(t, err)
	}

	// Act
	err = app.Restore(nil)
	if err != nil {
		require.NoError(t, err)
	}

	// Assert
	m := map[string][]string{
		"":                                 {"proto_deps.yml", "vendor.pb"},
		"/vendor.pb":                       {"github.com"},
		"/vendor.pb/github.com":            {"googleapis"},
		"/vendor.pb/github.com/googleapis": {"googleapis"},
		"/vendor.pb/github.com/googleapis/googleapis":            {"google"},
		"/vendor.pb/github.com/googleapis/googleapis/google":     {"api"},
		"/vendor.pb/github.com/googleapis/googleapis/google/api": {"annotations.proto"},
	}

	err = checkFolderContentEquality(t, testDir, m)

	if err != nil {
		log.Fatal(err)
		require.NoError(t, err)
	}

	defer os.RemoveAll(testDir)
}

func createdepsFileAndDir(t *testing.T, depsFileContent string) (string, error) {
	dirname := os.TempDir()
	testDir := path.Join(dirname, "proto_test", strconv.FormatInt(time.Now().Unix(), 10))
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		require.NoError(t, err)
	}

	create, err := os.Create(path.Join(testDir, "proto_deps.yml"))
	if err != nil {
		require.NoError(t, err)
	}

	_, err = create.WriteString(depsFileContent)
	if err != nil {
		require.NoError(t, err)
	}

	err = create.Close()
	if err != nil {
		require.NoError(t, err)
	}
	return testDir, err
}

func checkFolderContentEquality(t *testing.T, testDir string, m map[string][]string) error {
	err := filepath.Walk(testDir,
		func(currPath string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			curr := strings.Replace(currPath, testDir, "", 1)

			fmt.Println(curr)

			fileInfo, err := os.Stat(currPath)
			if err != nil {
				require.NoError(t, err)
			}

			if !fileInfo.IsDir() {
				return nil
			}

			entries, err := os.ReadDir(currPath)

			require.NoError(t, err)

			directoryEntities := funk.Map(entries, func(x os.DirEntry) string {
				return x.Name()
			})

			expectedFolderContext, ok := m[curr]

			if !ok {
				err = errors.New("unexpected item " + curr)
				require.NoError(t, err)
			}

			if !reflect.DeepEqual(directoryEntities, expectedFolderContext) {
				err = errors.New("invalid length")
				require.NoError(t, err)
			}

			return nil
		})

	return err
}
