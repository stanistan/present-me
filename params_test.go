package presentme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReviewParamsFromURL(t *testing.T) {
	params, err := ReviewParamsFromURL("stanistan/invoice-proxy/pull/3")
	require.NoError(t, err)
	require.Equal(t, int64(0), params.ReviewID)
}
