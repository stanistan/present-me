package github

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
)

const testDataDir = "testdata"

func TestDiffGenerator(t *testing.T) {
	testpaths, err := filepath.Glob(filepath.Join(testDataDir, "*.json"))
	assert.NoError(t, err)

	for _, testpath := range testpaths {
		_, filename := filepath.Split(testpath)
		test := filename[:len(filename)-len(filepath.Ext(testpath))]

		t.Run(test, func(t *testing.T) {
			var (
				data     = readFile(t, testpath)
				input    = diffFile(t, test, "input")
				expected = diffFile(t, test, "out")

				comment PullRequestComment
			)

			err = json.Unmarshal(data, &comment)
			assert.NoError(t, err, "test data file is an invalid json PullRequestComment")
			comment.DiffHunk = &input

			diff, err := generateDiff(&comment)
			assert.NoError(t, err, "failed to generate diff")
			assert.Equal(t, expected, diff+"\n", "generated diff doesn't match expected output")
		})
	}
}

func diffFile(t *testing.T, name, suffix string) string {
	t.Helper()
	bs := readFile(t, filepath.Join(testDataDir, name+"."+suffix+".diff"))
	return string(bs)
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	bs, err := os.ReadFile(path)
	assert.NoError(t, err, "could not read required test file")
	return bs
}
