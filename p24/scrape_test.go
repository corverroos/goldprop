package p24

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTile(t *testing.T) {
	tt := tile{
		Price: "R 1 050 000",
		Bedrooms: "0",
		Size: "80 m²",
	}
	p, err := tt.GetPrice()
	require.NoError(t, err)
	require.Equal(t, 1050000, p)

	b, err := tt.GetBedrooms()
	require.NoError(t, err)
	require.Equal(t, 0, b)

	s, err := tt.GetSize()
	require.NoError(t, err)
	require.Equal(t, 80, s)
}
