package parser

import (
	//	"ProtoDepsResolver/internal/models"
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
					`
{
  "version": "",
  "deps": [ ]
}`)
				mock.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)
			},
			err: true,
		},
		{
			name: "success_version",
			prepare: func(mock mock.MockIFileReader) {
				fileMock := []byte(
					`{
  "version": "v1",
  "deps": [ 
{
      "type": "path",
      "source": "/Users/Documents/ae/gitlab/platform/dotnet/main-service/src/Aer.RegressPlatform.Grpc/Api/vendor.pb/github.com/googleapis/googleapis/google/api/annotations.proto",
      "relativePath": "github.com/googleapis",
      "version": "v1"
    }]
}`)
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
			`
{
  "version": "v1",
  "deps": [
	{
      "type": "git",
      "source": "https://github.com/googleapis/googleapis.git",
      "relativePath": "google/api/http.proto",
      "tag": "common-protos-1_3_1"
    }
  ]
}
`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "common-protos-1_3_1", deps[0].Tag)
		require.Equal(t, "git", deps[0].Type)
		require.Equal(t, "https://github.com/googleapis/googleapis.git", deps[0].Source)
		require.Equal(t, "google/api/http.proto", deps[0].RelativePath)
		require.Equal(t, "", deps[0].CommitRevision)
		require.Equal(t, "", deps[0].Branch)
	})

	t.Run("TestGetDepsVersionTagWithRevision", func(t *testing.T) {
		fileMock := []byte(
			`
		{
			"version": "v1",
			"deps": [
			{
				"type": "git",
				"source": "https://github.com/googleapis/googleapis.git",
				"relativePath": "google/api/http.proto",
				"tag": "common-protos-1_3_1",
				"commitRevision":"409e134ffaac"
				}
			]
		}`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "", deps[0].Version)
		require.Equal(t, "https://github.com/googleapis/googleapis.git", deps[0].Source)
		require.Equal(t, "google/api/http.proto", deps[0].RelativePath)

		require.Equal(t, "git", deps[0].Type)
		require.Equal(t, "409e134ffaac", deps[0].CommitRevision)
	})

	t.Run("TestGetDepsUrl", func(t *testing.T) {
		fileMock := []byte(
			`
{
			"version": "v1",
			"deps": [
			{
				"type": "url",
				"source": "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto",
				"relativePath": "github.com/googleapis/googleapis/google/api/annotations.proto",
				"tag": "",
				"commitRevision":"",
				"version": "v1"
				}
			]
		}
`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1", deps[0].Version)
		require.Equal(t, "github.com/googleapis/googleapis/google/api/annotations.proto", deps[0].RelativePath)
		require.Equal(t, "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto", deps[0].Source)
		require.Equal(t, "", deps[0].CommitRevision)
		require.Equal(t, "", deps[0].Branch)
		require.Equal(t, "url", deps[0].Type)
	})

	t.Run("TestGetDepsPath", func(t *testing.T) {
		fileMock := []byte(
			`
{
  "version": "v1",
  "deps": [
    {
      "type": "path",
      "source": "/Users/src/greeter.proto",
      "relativePath": "github.com/googleapis",
      "version": "v1"
    }
  ]
}
`)
		fileReader.EXPECT().ReadFile(gomock.Any()).Return(fileMock, nil)

		parser := NewFileParser(fileReader)

		deps, err := parser.GetDeps("test")
		require.NoError(t, err)
		require.Len(t, deps, 1)
		require.Equal(t, "v1", deps[0].Version)
		require.Equal(t, "/Users/src/greeter.proto", deps[0].Source)
		require.Equal(t, "github.com/googleapis", deps[0].RelativePath)
		require.Equal(t, "path", deps[0].Type)
		require.Equal(t, "", deps[0].CommitRevision)
	})
}
