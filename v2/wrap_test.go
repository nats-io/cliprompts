package cliprompts

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrapSprintf(t *testing.T) {
	v := WrapSprintf(5, "%s", "01234 56789 abcdef")
	require.Equal(t, "01234\n56789\nabcdef", v)
}

func TestWrap(t *testing.T) {
	v := Wrap(5, "01234", "56789", "abcdef")
	require.Equal(t, "01234\n56789\nabcdef", v)
}
