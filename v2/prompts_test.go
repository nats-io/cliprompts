/*
 * Copyright 2018-2019 The NATS Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cliprompts

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailValidator(t *testing.T) {
	var tests = []struct {
		value      string
		shouldFail bool
	}{
		{"", false},
		{"a", true},
		{"@a", true},
		{"a@a", false}, // this is a bad, but valid email on a local network
		{"a@a.com", false},
	}

	o := NewEmailValidator()
	var opts Opts
	o(&opts)
	require.NotNil(t, opts.Fn)

	for _, vv := range tests {
		err := opts.Fn(vv.value)
		if err != nil && vv.shouldFail == false {
			t.Fatalf("expected not fail: %v", vv.value)
		} else if err == nil && vv.shouldFail {
			t.Fatalf("expected to fail on %v but didn't", vv.value)
		}
	}
}

func TestLengthValidator(t *testing.T) {
	o := NewLengthValidator(0)
	var opts Opts
	o(&opts)
	require.NotNil(t, opts.Fn)
	require.NoError(t, opts.Fn(""))
	require.NoError(t, opts.Fn("a"))
	require.NoError(t, opts.Fn("aaa"))

	o = NewLengthValidator(1)
	o(&opts)
	require.Error(t, opts.Fn(""))
	require.NoError(t, opts.Fn("a"))
	require.NoError(t, opts.Fn("aaa"))
}

func TestNewLengthValidator(t *testing.T) {
	o := NewLengthValidator(1)
	var opts Opts
	o(&opts)
	require.NotNil(t, opts.Fn)
}

func TestPathOrURLValidator(t *testing.T) {
	o := NewPathOrURLValidator()
	var opts Opts
	o(&opts)
	require.NotNil(t, opts.Fn)

	require.Error(t, opts.Fn(""))
	require.Error(t, opts.Fn("/tmp"))
	f, err := os.CreateTemp("", "file")
	require.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	require.NoError(t, opts.Fn(f.Name()))
	require.NoError(t, opts.Fn("http://www.google.com"))
}

func TestURLValidator(t *testing.T) {
	o := NewURLValidator("http", "https")
	var opts Opts
	o(&opts)
	require.NotNil(t, opts.Fn)

	require.Error(t, opts.Fn(""))
	require.Error(t, opts.Fn("/tmp"))
	require.Error(t, opts.Fn("https://"))
	require.Error(t, opts.Fn("ftp://localhost/path"))
	require.NoError(t, opts.Fn("http://localhost/path"))
	require.NoError(t, opts.Fn("http://localhost"))
	require.NoError(t, opts.Fn("HTTP://localhost"))
	require.NoError(t, opts.Fn("https://localhost"))
}

func TestItalic(t *testing.T) {
	v := Italic("test")
	require.Equal(t, fmt.Sprintf(italicTemplate, "test"), v)
}

func TestBold(t *testing.T) {
	v := Bold("test")
	require.Equal(t, fmt.Sprintf(boldTemplate, "test"), v)
}

func TestUnderline(t *testing.T) {
	v := Underline("test")
	require.Equal(t, fmt.Sprintf(underlineTemplate, "test"), v)
}

func TestHelp(t *testing.T) {
	fn := Help("hello")
	var opts Opts
	fn(&opts)
	require.Equal(t, "hello", opts.Help)
}

func TestVal(t *testing.T) {
	fn := Val(func(v string) error {
		return nil
	})
	var opts Opts
	fn(&opts)
	require.NotNil(t, opts.Fn)
}
