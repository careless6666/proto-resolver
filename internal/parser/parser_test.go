package parser

import (
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

//mock data

/* happy path
version: v1
deps:
  - git: github.com/googleapis/googleapis/google/api/http.proto v0.0.0-20211005231101-409e134ffaac
  - git: github.com/googleapis/googleapis/google/api/annotations.proto v0.0.0-20211005231101-409e134ffaac
  - url: https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto ./github.com/googleapis/googleapis/blob/master/google/api v1
  - path: /tmp/path ./github.com/googleapis v1
*/

/*
variations of format

deps:
  - github.com/googleapis/googleapis/google/api v0.0.0-20211005231101-409e134ffaac
  - github.com/googleapis/googleapis/google/api/* v0.0.0-20211005231101-409e134ffaac
  - github.com/googleapis/googleapis/google/api/* v1.0.0

*/
