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
	tests := []struct {
		name, in, out string
		typ           interface{}
	}{{
		// https://github.com/yaml/yaml/issues/48
		name: "merging a sequence",
		in: `array1: &my_array_alias
- foo
- bar

array2:
- <<: *my_array_alias
- zap`,
		typ: make(map[string][]string),
		out: "array1:\n    - foo\n    - bar\narray2:\n    - foo\n    - bar\n    - zap\n",
	}, {
		// https://github.com/yaml/yaml/issues/48
		name: "merging two sequences",
		in: `array1: &my_array_alias
- foo
- bar

array2:
- first
- <<: [*my_array_alias, *my_array_alias]
- zap`,
		typ: make(map[string][]string),
		out: "array1:\n    - foo\n    - bar\narray2:\n    - first\n    - foo\n    - bar\n    - foo\n    - bar\n    - zap\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, xyaml.Unmarshal([]byte(tt.in), tt.typ))
			newData, err := xyaml.Marshal(tt.typ)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, string(newData))
		})
	}

	// todo test: sequences that contain arrays with a merge
	// <<: *my_array_alias, *my_array_alias
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
