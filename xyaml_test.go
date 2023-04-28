/*
 * SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
 *
 * SPDX-License-Identifier: MIT
 */

package xyaml_test

import (
	"math"
	"reflect"
	"testing"

	"codeberg.org/6543/xyaml"
	"github.com/stretchr/testify/assert"
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
	}}
	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			newValue := reflect.New(reflect.ValueOf(tt.value).Type())
			xyaml.Unmarshal([]byte(tt.data), newValue.Interface())
			assert.EqualValues(t, tt.value, newValue.Elem().Interface())
		})
	}
}
