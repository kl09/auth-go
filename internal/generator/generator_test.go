package generator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kl09/auth-go/internal/generator"
)

func TestGenerateRandomString(t *testing.T) {
	s, err := generator.GenerateRandomString(128)
	require.Nil(t, err)
	require.Len(t, s, 128)

	s2, err := generator.GenerateRandomString(128)
	require.Nil(t, err)
	require.Len(t, s2, 128)
	require.NotEqual(t, s, s2)

	s3, err := generator.GenerateRandomString(1)
	require.Nil(t, err)
	require.Len(t, s3, 1)
}
