package datasizes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/images/pkg/datasizes"
)

func TestDataSizeToUint64(t *testing.T) {
	cases := []struct {
		input   string
		success bool
		output  uint64
	}{
		{"123", true, 123},
		{"123 kB", true, 123000},
		{"123 KiB", true, 123 * 1024},
		{"123 MB", true, 123 * 1000 * 1000},
		{"123 MiB", true, 123 * 1024 * 1024},
		{"123 GB", true, 123 * 1000 * 1000 * 1000},
		{"123 GiB", true, 123 * 1024 * 1024 * 1024},
		{"123 TB", true, 123 * 1000 * 1000 * 1000 * 1000},
		{"123 TiB", true, 123 * 1024 * 1024 * 1024 * 1024},
		{"123kB", true, 123000},
		{"123KiB", true, 123 * 1024},
		{" 123  ", true, 123},
		{"  123kB  ", true, 123000},
		{"  123KiB  ", true, 123 * 1024},
		{"123 KB", false, 0},
		{"123 mb", false, 0},
		{"123 PB", false, 0},
		{"123 PiB", false, 0},
	}

	for _, c := range cases {
		result, err := datasizes.Parse(c.input)
		if c.success {
			require.Nil(t, err)
			assert.EqualValues(t, c.output, result)
		} else {
			assert.NotNil(t, err)
		}
	}
}
