// SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
// SPDX-License-Identifier: MIT

package xyaml

import (
	"gopkg.in/yaml.v3"
)

// Source: https://github.com/go-yaml/yaml/blob/3e3283e801afc229479d5fc68aa41df1137b8394/resolve.go#L80
const mergeTag = "!!merge"

// Unmarshal use underlying yaml.v3 lib func to decode and alter based on the extensions afterwards
func Unmarshal(in []byte, out interface{}) (err error) {
	// internal unmarshal
	node := new(yaml.Node)
	if err := yaml.Unmarshal(in, node); err != nil {
		return err
	}

	// process extensions
	if err := mergeSequences(node, 0); err != nil {
		return err
	}

	// unmarshal to final dest type
	return node.Decode(out)
}

// Marshal just use the normal underlying yaml.v3 lib func
func Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}
