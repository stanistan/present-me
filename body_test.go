package crap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrippingStartingNumber(t *testing.T) {
	data := "1. First Things First"
	require.Equal(t, "First Things First", stripLeadingNumber(data))
}
