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
	"errors"
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Logger func(args ...interface{})

type Validator func(string) error

var cli PromptLib

// set a Logger during a test (cli.LogFn = t.Log) to debug interactive prompts
var LogFn Logger

// lint:ignore U1000
var output io.Writer = os.Stdout

type PromptLib interface {
	Prompt(label string, value string, opt ...Opt) (string, error)
	Confirm(m string, defaultValue bool, o ...Opt) (bool, error)
	Password(m string, o ...Opt) (string, error)
	Select(m string, value string, choices []string, o ...Opt) (int, error)
	MultiSelect(m string, choices []string, o ...Opt) ([]int, error)
}

type Opts struct {
	Help string
	Fn   Validator
}

type Opt func(o Opts)

func processOpts(opt ...Opt) Opts {
	var opts Opts
	for _, o := range opt {
		o(opts)
	}
	return opts
}

func Help(h string) Opt {
	return func(o Opts) {
		o.Help = h
	}
}

func Val(fn Validator) Opt {
	return func(o Opts) {
		o.Fn = fn
	}
}

func init() {
	ResetPromptLib()
	LogFn = nil
}

func SetPromptLib(p PromptLib) {
	cli = p
}

func ResetPromptLib() {
	SetPromptLib(&SurveyUI{})
	LogFn = nil
}

func SetOutput(out io.Writer) {
	output = out
}

func GetOutput() io.Writer {
	return output
}

func Underline(s string) string {
	return fmt.Sprintf("\xff\033[4m\xff%s\xff\033[0m\xff", s)
}

func Bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func Italic(s string) string {
	return fmt.Sprintf("\033[3m%s\033[0m", s)
}

func Prompt(label string, value string, o ...Opt) (string, error) {
	return cli.Prompt(label, value, o...)
}

func ConfirmYN(m string, o ...Opt) (bool, error) {
	return cli.Confirm(m, true, o...)
}

func ConfirmNY(m string, o ...Opt) (bool, error) {
	return cli.Confirm(m, false, o...)
}

func Password(m string, o ...Opt) (string, error) {
	return cli.Password(m, o...)
}

func Select(m string, value string, choices []string, o ...Opt) (int, error) {
	return cli.Select(m, value, choices, o...)
}

func MultiSelect(m string, choices []string, o ...Opt) ([]int, error) {
	return cli.MultiSelect(m, choices, o...)
}

func NewEmailValidator() Opt {
	return Val(EmailValidator())
}

func EmailValidator() Validator {
	return func(input string) error {
		if input != "" {
			_, err := mail.ParseAddress(input)
			return err
		}
		return nil
	}
}

func NewLengthValidator(min int) Opt {
	return Val(LengthValidator(min))
}

func LengthValidator(min int) Validator {
	return func(input string) error {
		if len(input) >= min {
			return nil
		}
		return errors.New("value is too short")
	}
}

func NewPathOrURLValidator() Opt {
	return Val(PathOrURLValidator())
}

func PathOrURLValidator() Validator {
	return func(s string) error {
		if u, err := url.Parse(s); err == nil && u.Scheme != "" {
			return nil
		}

		v, err := homedir.Expand(s)
		if err != nil {
			return err
		}
		v, err = filepath.Abs(v)
		if err != nil {
			return err
		}
		info, err := os.Stat(v)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return errors.New("path is not a file")
		}
		return nil
	}
}

func NewURLValidator(protocol ...string) Opt {
	return Val(URLValidator(protocol...))
}

func URLValidator(protocol ...string) Validator {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errors.New("url cannot be empty")
		}
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		scheme := strings.ToLower(u.Scheme)

		ok := false
		for _, v := range protocol {
			if scheme == v {
				ok = true
				break
			}
		}
		if !ok {
			var protos []string
			protos = append(protos, protocol...)
			return fmt.Errorf("scheme %q is not supported (%v)", scheme, strings.Join(protos, ", "))
		}
		if u.Host == "" {
			return fmt.Errorf("no host specified (%v)", s)
		}

		return nil
	}
}
