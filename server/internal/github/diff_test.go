package github

import (
	_ "embed"
	"testing"

	"github.com/stanistan/present-me/internal/github/diff"
	"github.com/stretchr/testify/require"
)

func TestParseDiffHunkPrefix(t *testing.T) {
	t.Run("first case", func(t *testing.T) {
		meta, err := diff.ParseHunkMeta("@@ -230,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.NoError(t, err)
		require.Equal(t, diff.HunkMeta{
			Original: diff.HunkRange{230, 6, "+"},
			New:      diff.HunkRange{200, 9, "-"},
		}, meta)
	})
	t.Run("succeeds", func(t *testing.T) {
		meta, err := diff.ParseHunkMeta("@@ -0,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.NoError(t, err)
		require.Equal(t, diff.HunkMeta{
			Original: diff.HunkRange{0, 6, "+"},
			New:      diff.HunkRange{200, 9, "-"},
		}, meta)
	})
	t.Run("errs", func(t *testing.T) {
		_, err := diff.ParseHunkMeta("@@ -0,6 +a,9 @@ if (!defined $initial_reply_to && $prompting) {")
		require.Error(t, err)
	})
}

//go:embed test_hunk.diff
var diffHunk string

func TestParser(t *testing.T) {
	right := "RIGHT"
	line := 185
	comment := &PullRequestComment{DiffHunk: &diffHunk, Side: &right, Line: &line}
	diff, err := generateDiff(comment)
	require.NoError(t, err)
	require.Equal(t, ` type serviceContext struct {
-	Management ManagementType
-	Actions    *util.StringSet
+	Actions *util.StringSet
 }`, diff)
}
