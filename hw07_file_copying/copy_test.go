package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inputPath = "./testdata/input.txt"
	outPath   = "./testdata/out.txt"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
		limit  int64
	}{
		{name: "Offset and limit are 0", offset: 0, limit: 0},
		{name: "Offset is 0 and limit is less than file size", offset: 0, limit: 10},
		{name: "Offset is 0 and limit is less than file size", offset: 0, limit: 1000},
		{name: "Offset is 0 and limit is bigger than file size", offset: 0, limit: 10000},
		{name: "Offset is less than limit and bigger than 0", offset: 100, limit: 1000},
		{name: "Offset bigger than limit", offset: 6000, limit: 1000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(inputPath, outPath, tc.offset, tc.limit, nil)
			expectedPath := fmt.Sprintf("./testdata/out_offset%d_limit%d.txt", tc.offset, tc.limit)
			expected, _ := os.ReadFile(expectedPath)
			out, _ := os.ReadFile(outPath)

			require.Equal(t, expected, out)
			require.NoError(t, err)
		})
	}

	t.Run("Offset is bigger than file size", func(t *testing.T) {
		err := Copy(inputPath, outPath, 10000, 0, nil)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	if err := os.Remove(outPath); err != nil {
		panic("Couldn't delete output file")
	}
}

func TestCopyLoadig(t *testing.T) {
	loadingCh := make(chan int64)
	go func() {
		err := Copy(inputPath, outPath, 0, 100, loadingCh)
		require.NoError(t, err)
	}()

	i := int64(1)
	for c := range loadingCh {
		require.Equal(t, i, c)
		i++
	}

	if err := os.Remove(outPath); err != nil {
		panic("Couldn't delete output file")
	}
}
