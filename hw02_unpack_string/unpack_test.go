package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	testStr := `test
2
er3

 2`

	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "  2  ", expected: "     "},
		// non-printed symbols and escape symbols in string
		{input: "aa2\n2\b3", expected: "aaa\n\n\b\b\b"},
		{input: "aa2\\n2\\t3", expected: "aaa\\nn\\ttt"},
		// encoded strings
		{input: "\u65e53\u672c2\u8a9e", expected: "\u65e5\u65e5\u65e5\u672c\u672c\u8a9e"},
		{input: "\U000065e52\U0000672c\U00008a9e2", expected: "\U000065e5\U000065e5\U0000672c\U00008a9e\U00008a9e"},
		// backquoted strings
		{input: `a2\n3\2a`, expected: `aa\nnn\\a`},
		{input: testStr, expected: "test\n\n\nerrr\n\n  "},

		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "aa2tt34", "a3456T"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
