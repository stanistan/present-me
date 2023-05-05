package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDiffHunkPrefix(t *testing.T) {
	t.Run("first case", func(t *testing.T) {
		meta, err := diffHunkPrefix("@@ -230,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.NoError(t, err)
		require.Equal(t, diffHunkMetadata{230, 200}, meta)
	})
	t.Run("succeeds", func(t *testing.T) {
		meta, err := diffHunkPrefix("@@ -0,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.NoError(t, err)
		require.Equal(t, diffHunkMetadata{0, 200}, meta)
	})
	t.Run("errs", func(t *testing.T) {
		_, err := diffHunkPrefix("@@ -0,6 +a,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.Error(t, err)
	})
}
