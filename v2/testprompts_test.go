package cliprompts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeInputs() []interface{} {
	var a []interface{}
	a = append(a, "a", false, true, "secret", 1, []int{1, 2})
	return a
}

func Test_PromptTypeCheck(t *testing.T) {
	cli = NewTestPrompts(makeInputs())
	v, err := Prompt("m", "b")
	require.NoError(t, err)
	require.Equal(t, "a", v)

	fv, err := Confirm("m", true)
	require.NoError(t, err)
	require.False(t, fv)

	tv, err := Confirm("m", false)
	require.NoError(t, err)
	require.True(t, tv)

	sv, err := Password("m")
	require.NoError(t, err)
	require.Equal(t, "secret", sv)

	cv, err := Select("m", "a", []string{"a", "b", "c"})
	require.NoError(t, err)
	require.Equal(t, 1, cv)

	av, err := MultiSelect("m", []string{"a", "b", "c"})
	require.NoError(t, err)
	require.Len(t, av, 2)
	require.Equal(t, 1, av[0])
	require.Equal(t, 2, av[1])
}

func Test_ErrIfBadTestInput(t *testing.T) {
	cli = NewTestPrompts([]interface{}{"a"})
	_, err := Confirm("m", true)
	require.EqualError(t, err, "m confirm expected a bool: a")

	cli = NewTestPrompts([]interface{}{"a"})
	_, err = Confirm("m", false)
	require.EqualError(t, err, "m confirm expected a bool: a")

	cli = NewTestPrompts([]interface{}{true})
	_, err = Password("m")
	require.EqualError(t, err, "m password expected a string: true")

	cli = NewTestPrompts([]interface{}{"a"})
	_, err = Select("m", "a", []string{"a", "b", "c"})
	require.EqualError(t, err, "m select expected an int: a")

	cli = NewTestPrompts([]interface{}{"a"})
	_, err = MultiSelect("m", []string{"a", "b", "c"})
	require.EqualError(t, err, "m multiselect expected []int: a")
}

func Test_TestPromptsValidator(t *testing.T) {
	cli := NewTestPrompts([]interface{}{"a"})
	var called bool
	v, err := cli.Prompt("test", "x", Val(func(s string) error {
		called = true
		return nil
	}))

	require.NoError(t, err)
	require.Equal(t, "a", v)
	require.True(t, called)
}

func Test_MultiSelectNoValidator(t *testing.T) {
	_, err := MultiSelect("test", []string{"a", "b"}, Val(func(s string) error {
		return nil
	}))
	require.Error(t, err)
	require.Equal(t, errValidatorNotSupported, err)
}

func Test_SelectNoValidator(t *testing.T) {
	_, err := Select("test", "a", []string{"a", "b"}, Val(func(s string) error {
		return nil
	}))
	require.Error(t, err)
	require.Equal(t, errValidatorNotSupported, err)
}

func Test_ConfirmValidator(t *testing.T) {
	_, err := Confirm("test", true, Val(func(s string) error {
		return nil
	}))
	require.Error(t, err)
	require.Equal(t, errValidatorNotSupported, err)
}
