package downloader

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRepoName(t *testing.T) {
	// TODO: add more tests
	repoName, err := GetRepoName("https://github.com/googleapis/googleapis.git")
	require.Equal(t, nil, err)
	require.Equal(t, "googleapis", repoName)
}
