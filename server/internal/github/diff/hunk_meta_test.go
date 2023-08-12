package diff

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestParseHunkMeta(t *testing.T) {
	t.Run("first case", func(t *testing.T) {
		meta, err := ParseHunkMeta("@@ -230,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		assert.NoError(t, err)
		assert.Equal(t, HunkMeta{
			Original: HunkRange{StartingAt: 230, NumLines: 6, IgnorePrefix: "+"},
			New:      HunkRange{StartingAt: 200, NumLines: 9, IgnorePrefix: "-"},
		}, meta)
	})
	t.Run("succeeds", func(t *testing.T) {
		meta, err := ParseHunkMeta("@@ -0,6 +200,9 @@ if (!defined $initial_reply_to && $prompting) {")
		assert.NoError(t, err)
		assert.Equal(t, HunkMeta{
			Original: HunkRange{StartingAt: 0, NumLines: 6, IgnorePrefix: "+"},
			New:      HunkRange{StartingAt: 200, NumLines: 9, IgnorePrefix: "-"},
		}, meta)
	})
	t.Run("errs", func(t *testing.T) {
		_, err := ParseHunkMeta("@@ -0,6 +a,9 @@ if (!defined $initial_reply_to && $prompting) {")
		assert.Error(t, err)
	})
}
