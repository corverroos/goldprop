package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchUrl(t *testing.T) {
	url, err := makeSearchUrl(0,"Sea Point", "Steenberg")
	require.NoError(t, err)
	require.Equal(t,
		"https://www.property24.com/apartments-for-sale/advanced-search/results?sp=s%3D11021%2C9040",
		url)

	url, err = makeSearchUrl(1000000,"Retreat Industrial")
	require.NoError(t, err)
	require.Equal(t,
		"https://www.property24.com/apartments-for-sale/advanced-search/results?pf=700000&sp=s%3D25704",
		url)
}

