/*
 * SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
 *
 * SPDX-License-Identifier: MIT
 */

package xyaml_test

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"codeberg.org/6543/xyaml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestStdDecode(t *testing.T) {
	tests := []struct {
		data  string
		value interface{}
	}{{
		"value: hi",
		map[string]string{"value": "hi"},
	}, {
		"boolean: TRUE",
		map[string]interface{}{"boolean": true},
	}, {
		"a: {b: c}",
		&struct{ A *map[string]string }{&map[string]string{"b": "c"}},
	}, {
		"float32_max: 3.40282346638528859811704183484516925440e+38",
		map[string]float32{"float32_max": math.MaxFloat32},
	}, {
		"a: [1, 2]",
		&struct{ A [2]int }{[2]int{1, 2}},
	}, {
		"a: {b: c}",
		&struct{ A *struct{ B string } }{&struct{ B string }{"c"}},
	}, {
		"a:\n - 1\n - &ref 2\nb:\n - *ref\n - 3",
		map[string][]string{"a": {"1", "2"}, "b": {"2", "3"}},
	}}
	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			newValue := reflect.New(reflect.ValueOf(tt.value).Type())
			xyaml.Unmarshal([]byte(tt.data), newValue.Interface())
			assert.EqualValues(t, tt.value, newValue.Elem().Interface())
		})
	}
}

func TestMergeSequences(t *testing.T) {
	// https://github.com/yaml/yaml/issues/48
	tests := []struct {
		name, in, out string
	}{{
		name: "merging a sequence",
		in: `array1: &my_array_alias
- foo
- bar

array2:
- <<: *my_array_alias
- zap`,
		out: "array1:\n    - foo\n    - bar\narray2:\n    - foo\n    - bar\n    - zap\n",
	}, {
		name: "merging two sequences",
		in: `array1: &my_array_alias
- foo
- bar

array2:
- first
- <<: [*my_array_alias, *my_array_alias]
- zap`,
		out: "array1:\n    - foo\n    - bar\narray2:\n    - first\n    - foo\n    - bar\n    - foo\n    - bar\n    - zap\n",
	}, {
		name: "merge sequences independent",
		in: `array1: &alias1
- one
- two
result:
- 1
- <<: *alias1
- 2
- <<: *alias1
`,
		out: "array1:\n    - one\n    - two\nresult:\n    - \"1\"\n    - one\n    - two\n    - \"2\"\n    - one\n    - two\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := make(map[string][]string)
			assert.NoError(t, xyaml.Unmarshal([]byte(tt.in), out))
			newData, err := xyaml.Marshal(out)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, string(newData))
		})
	}
}

func TestMergeSequenceErrors(t *testing.T) {
	tests := []struct {
		name, in, errStg string
		err              error
	}{{
		name: "test",
		in: `array2:
- <<: a, b
- zap`,
		errStg: "",
		err:    xyaml.ErrSequenceMerge,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := new(yaml.Node)
			err := xyaml.Unmarshal([]byte(tt.in), node)
			if assert.Error(t, err) {
				if assert.Truef(t, errors.Is(err, tt.err), "want: '%s' error, but got '%s' error", tt.err, err) && tt.errStg != "" {
					assert.EqualValues(t, tt.errStg, err.Error())
				}
			}
		})
	}
}

func TestMergeMap(t *testing.T) {
	tests := []struct {
		name, in, out string
	}{{
		name: "merge unique maps",
		in:   "letters: &letters\n  a: \"A\"\nnumbers: &numbers\n  one: \"1\"\ncombined:\n  <<: [ *letters, *numbers ]",
		out:  "combined:\n    a: A\n    one: \"1\"\nletters:\n    a: A\nnumbers:\n    one: \"1\"\n",
	}, {
		name: "extend maps",
		in:   "base: &base\n  foo: FOO\n  bar: BAR\nextended:\n  <<: *base\n  zap: ZAP",
		out:  "base:\n    bar: BAR\n    foo: FOO\nextended:\n    bar: BAR\n    foo: FOO\n    zap: ZAP\n",
	}, {
		name: "overwrite vales",
		in:   "\nbase: &base\n  val: ONE\n  next: 1\n  zap: 3\noverwrite:\n  next: 2\n  <<: *base\n  val: TWO\n",
		out:  "base:\n    next: \"1\"\n    val: ONE\n    zap: \"3\"\noverwrite:\n    next: \"2\"\n    val: TWO\n    zap: \"3\"\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := make(map[string]map[string]string)
			assert.NoError(t, yaml.Unmarshal([]byte(tt.in), out))
			newData, err := xyaml.Marshal(out)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, string(newData))
		})
	}
}
