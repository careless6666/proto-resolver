package parser

import (
	"ProtoDepsResolver/internal/models"
	"ProtoDepsResolver/internal/parser/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestParserVersion(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(mock mock.MockIFileReader)
		err     bool
	}{
		{
			name: "error_version",
			prepare: func(mock mock.MockIFileReader) {
				fileMock := []byte(
					`verdsion: v1
deps:
  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac`)
				mock.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)
			},
			err: true,
		},
		{
			name: "success_version",
			prepare: func(mock mock.MockIFileReader) {
				fileMock := []byte(
					`version: v1
deps:
  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac`)
				mock.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := gomock.NewController(t)
			fileReader := mock.NewMockIFileReader(gc)

			tt.prepare(*fileReader)

			parser := NewFileParser(fileReader)

			_, err := parser.GetDeps("test")

			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetDepsVersion(t *testing.T) {
	gc := gomock.NewController(t)
	fileReader := mock.NewMockIFileReader(gc)

	t.Run("TestGetDepsVersionTagOnly", func(t *testing.T) {
		fileMock := []byte(
			`version: v1
deps:
  - git: https://github.com/googleapis/googleapis.git /google/api/http.proto v1.0.0`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1.0.0", deps[0].Version.Tag)
		require.Equal(t, "https://github.com/googleapis/googleapis.git", deps[0].Path)
		require.Equal(t, "/google/api/http.proto", deps[0].GitPath)
		require.Equal(t, "", deps[0].DestinationPath)
		require.Equal(t, models.DependencyTypeGit, deps[0].Type)
	})

	t.Run("TestGetDepsVersionTagWithRevision", func(t *testing.T) {
		fileMock := []byte(
			`version: v1
deps:
  - git: https://github.com/googleapis/googleapis.git /google/api/http.proto v1.0.0-20211005231101-409e134ffaac`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1.0.0", deps[0].Version.Tag)
		require.Equal(t, "https://github.com/googleapis/googleapis.git", deps[0].Path)
		require.Equal(t, "/google/api/http.proto", deps[0].GitPath)
		require.Equal(t, "", deps[0].DestinationPath)
		require.Equal(t, models.DependencyTypeGit, deps[0].Type)
		require.Equal(t, "409e134ffaac", deps[0].Version.CommitRevision)
	})

	t.Run("TestGetDepsUrl", func(t *testing.T) {
		fileMock := []byte(
			`version: v1
deps:
  - url: https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto ./github.com/googleapis/googleapis/blob/master/google/api v1`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1", deps[0].Version.Tag)
		require.Equal(t, "https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto", deps[0].Path)
		require.Equal(t, "./github.com/googleapis/googleapis/blob/master/google/api", deps[0].DestinationPath)
		require.Equal(t, models.DependencyTypeURL, deps[0].Type)
		require.Equal(t, "", deps[0].Version.CommitRevision)
	})

	t.Run("TestGetDepsPath", func(t *testing.T) {
		fileMock := []byte(
			`version: v1
deps:
  - path: /tmp/path ./github.com/googleapis v1`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1", deps[0].Version.Tag)
		require.Equal(t, "/tmp/path", deps[0].Path)
		require.Equal(t, "./github.com/googleapis", deps[0].DestinationPath)
		require.Equal(t, models.DependencyTypePath, deps[0].Type)
		require.Equal(t, "", deps[0].Version.CommitRevision)
	})
}

//mock data

/* happy path
version: v1
deps:
#  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac
  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac
  - git: github.com/googleapis/googleapis/google/api/annotations.proto v0.0.0-20211005231101-409e134ffaac
  - url: https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto ./github.com/googleapis/googleapis/blob/master/google/api v1
  - path: /tmp/path ./github.com/googleapis v1
*/

/*
variations of format

deps:
  - github.com/googleapis/googleapis/google/api v0.0.0-20211005231101-409e134ffaac
  - github.com/googleapis/googleapis/google/api v1.0.0

*/
