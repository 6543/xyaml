/*
 * SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
 *
 * SPDX-License-Identifier: MIT
 */

package xyaml

import (
	"gopkg.in/yaml.v3"
)

// Unmarshal use underlying yaml.v3 lib func to decode and alter based on the extensions afterwards
func Unmarshal(in []byte, out interface{}) (err error) {
	node := new(yaml.Node)
	yaml.Unmarshal(in, node)

	// TODO: do add extensions

	return node.Decode(out)
}

// Marshal just use the normal underlying yaml.v3 lib func
func Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}
