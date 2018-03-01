package storage

import (
	"fmt"
	"testing"
)

func TestIsErrorKeytNotFound(t *testing.T) {
	type Case struct {
		msg    string
		err    error
		expect bool
	}

	cases := []Case{
		{
			"KeyNotFoundError()",
			KeyNotFoundError([]byte("an_key")),
			true,
		},
		{
			"ErrorKeyNotFound{} struct",
			ErrorKeyNotFound{fmt.Errorf("struct")},
			true,
		},
		{
			"*ErrorKeyNotFound{} pointer",
			&ErrorKeyNotFound{fmt.Errorf("pointer")},
			true,
		},
		{
			"InternalError()",
			InternalError("go wrong"),
			false,
		},
		{
			"usual error",
			fmt.Errorf("usual error"),
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			actual := IsErrorKeyNotFound(c.err)
			if actual != c.expect {
				t.Fatalf("(%v) want %v got %v", c.err, c.expect, actual)
			}
		})
	}
}
