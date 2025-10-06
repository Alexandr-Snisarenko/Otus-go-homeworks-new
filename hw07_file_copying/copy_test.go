package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		desc      string
		src       string
		offset    int64
		limit     int64
		rewrite   bool
		outLength int64
		expectErr error
	}{
		{
			desc:      "offset > file length",
			src:       "1234567890",
			offset:    100,
			limit:     0,
			rewrite:   true,
			outLength: 0,
			expectErr: ErrOffsetExceedsFileSize,
		},
		{
			desc:      "offset",
			src:       "1234567890",
			offset:    5,
			limit:     0,
			rewrite:   true,
			outLength: 5,
			expectErr: nil,
		},
		{
			desc:      "offset+limit",
			src:       "1234567890",
			offset:    5,
			limit:     3,
			rewrite:   true,
			outLength: 3,
			expectErr: nil,
		},
		{
			desc:      "limit > length",
			src:       "1234567890",
			offset:    5,
			limit:     30,
			rewrite:   true,
			outLength: 5,
			expectErr: nil,
		},
		{
			desc:      "limit == length",
			src:       "1234567890",
			offset:    0,
			limit:     10,
			rewrite:   true,
			outLength: 10,
			expectErr: nil,
		},
		{
			desc:      "offset == length",
			src:       "1234567890",
			offset:    10,
			limit:     0,
			rewrite:   true,
			outLength: 0,
			expectErr: nil,
		},
		{
			desc:      "1 bite",
			src:       "1234567890",
			offset:    0,
			limit:     1,
			rewrite:   true,
			outLength: 1,
			expectErr: nil,
		},
		{
			desc:      "empty file",
			src:       "",
			offset:    0,
			limit:     0,
			rewrite:   true,
			outLength: 0,
			expectErr: nil,
		},
		{
			desc:      "dest file exists",
			src:       "1234567890",
			offset:    0,
			limit:     0,
			rewrite:   false,
			outLength: 0,
			expectErr: ErrFileExists,
		},
		{
			desc:      "unsuported in_file (dir)",
			src:       "",
			offset:    0,
			limit:     0,
			rewrite:   true,
			outLength: 0,
			expectErr: ErrUnsupportedFile,
		},
		{
			desc:      "unsuported out_file (dir)",
			src:       "1234567890",
			offset:    0,
			limit:     0,
			rewrite:   true,
			outLength: 0,
			expectErr: ErrUnsupportedFile,
		},
	}

	for _, tc := range tests {
		var err error
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			fIn, er := os.CreateTemp("", "hw07.test.*.txt")
			if er != nil {
				panic("can't create in test file")
			}
			defer fIn.Close()
			defer os.Remove(fIn.Name())

			fIn.Write([]byte(tc.src))

			switch tc.desc {
			case "dest file exists":
				err = Copy(fIn.Name(), fIn.Name(), tc.offset, tc.limit, tc.rewrite)
			case "unsuported in_file (dir)":
				dir, _ := filepath.Split(fIn.Name())
				err = Copy(dir, fIn.Name()+".out", tc.offset, tc.limit, tc.rewrite)
			case "unsuported out_file (dir)":
				dir, _ := filepath.Split(fIn.Name())
				err = Copy(fIn.Name(), dir, tc.offset, tc.limit, tc.rewrite)
			default:
				err = Copy(fIn.Name(), fIn.Name()+".out", tc.offset, tc.limit, tc.rewrite)
			}

			if tc.expectErr != nil {
				require.True(t, errors.Is(err, tc.expectErr))
			} else {
				defer os.Remove(fIn.Name() + ".out")
				require.FileExists(t, fIn.Name()+".out")
				fInfo, _ := os.Stat(fIn.Name() + ".out")
				require.Equal(t, tc.outLength, fInfo.Size())
			}
		})
	}
}
