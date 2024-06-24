package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRepoPathFromAddress(t *testing.T) {
	// TODO: add more tests
	repoName, err := GetRepoPathFromAddress("https://github.com/googleapis/googleapis.git")
	require.Equal(t, nil, err)
	require.Equal(t, "github.com/googleapis/googleapis", repoName)
}
