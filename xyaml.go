// SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
// SPDX-License-Identifier: MIT

package xyaml

import (
	"gopkg.in/yaml.v3"
)

type Parser interface {
	Unmarshal(in []byte, out any) (err error)
	Marshal(in any) (out []byte, err error)
	MergeSequences(in *yaml.Node) (err error)
}

// DefaultParser uses default config,
// can be overloaded to set new config for public functions.
var DefaultParser = NewParser()

// Unmarshal use underlying yaml.v3 lib func to decode and alter based on the extensions afterwards with default config
func Unmarshal(in []byte, out any) (err error) {
	return DefaultParser.Unmarshal(in, out)
}

// Marshal just use the normal underlying yaml.v3 lib func
func Marshal(in any) (out []byte, err error) {
	return DefaultParser.Marshal(in)
}

// MergeSequences recursively search for sequence merge indicator "<<:" and merge if detected
func MergeSequences(node *yaml.Node) error {
	return DefaultParser.MergeSequences(node)
}
