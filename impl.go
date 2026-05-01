// SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
// SPDX-License-Identifier: MIT

package xyaml

import "gopkg.in/yaml.v3"

type config struct {
	maxDepth uint16
}

// MergeSequences recursively search for sequence merge indicator "<<:" and merge if detected
func (c config) MergeSequences(node *yaml.Node) error {
	return defaultConfig.mergeSequences(node, 0)
}

// Unmarshal use underlying yaml.v3 lib func to decode and alter based on the extensions afterwards
func (c config) Unmarshal(in []byte, out any) (err error) {
	// internal unmarshal
	node := new(yaml.Node)
	if err := yaml.Unmarshal(in, node); err != nil {
		return err
	}

	// process extensions
	if err := c.mergeSequences(node, 0); err != nil {
		return err
	}

	// unmarshal to final dest type
	return node.Decode(out)
}

// Marshal just use the normal underlying yaml.v3 lib func
func (c config) Marshal(in any) (out []byte, err error) {
	return yaml.Marshal(in)
}
