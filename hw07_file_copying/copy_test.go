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
			err := Copy(inputPath, outPath, tc.offset, tc.limit)
			if err != nil {
				require.NoError(t, err, "Error copying file")
			}

			expectedPath := fmt.Sprintf("./testdata/out_offset%d_limit%d.txt", tc.offset, tc.limit)
			expected, err := os.ReadFile(expectedPath)
			if err != nil {
				require.NoError(t, err, "Error reading expected file")
			}

			out, _ := os.ReadFile(outPath)
			if err != nil {
				require.NoError(t, err, "Error reading output file")
			}

			require.Equal(t, expected, out)
		})
	}

	t.Run("Offset is bigger than file size", func(t *testing.T) {
		err := Copy(inputPath, outPath, 10000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	if err := os.Remove(outPath); err != nil {
		require.NoError(t, err, "Error removing output file")
	}
}

func TestCopyErrors(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
		limit  int64
		err    error
	}{
		{name: "Offset is bigger than file size", offset: 10000, limit: 0, err: ErrOffsetExceedsFileSize},
		{name: "Offset is negative", offset: -10, limit: 0, err: ErrInvalidOffset},
		{name: "Limit is negative", offset: 0, limit: -10, err: ErrLimitOffset},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(inputPath, outPath, tc.offset, tc.limit)
			require.ErrorIs(t, tc.err, err)
		})
	}
}

func TestCopySameFiles(t *testing.T) {
	t.Run("From file and to file are the same", func(t *testing.T) {
		path := "./testdata/same_file_input.txt"
		err := os.WriteFile(path, []byte("From file and to file are the same"), 0o644)
		if err != nil {
			require.NoError(t, err, "Error writing file")
		}

		err = Copy(path, path, 0, 10)
		if err != nil {
			require.NoError(t, err, "Error copying file")
		}

		file, err := os.ReadFile(path)
		if err != nil {
			require.NoError(t, err, "Error reading file")
		}

		require.Equal(t, "From file ", string(file))
	})
}
