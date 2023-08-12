package github

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestReviewParamsFromURL(t *testing.T) {
	params, err := ReviewParamsFromURL("stanistan/invoice-proxy/pull/3")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), params.ReviewID)
}
