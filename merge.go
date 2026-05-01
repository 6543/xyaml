// SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
// SPDX-License-Identifier: MIT

package xyaml

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Source: https://github.com/go-yaml/yaml/blob/3e3283e801afc229479d5fc68aa41df1137b8394/resolve.go#L80
const mergeTag = "!!merge"

var (
	// ErrMaxDepth indicates there is likely a loop that got caught
	ErrMaxDepth = errors.New("max depth reached")
	// ErrBrokenMappingNode indicates a broken map node
	ErrBrokenMappingNode = errors.New("broken mapping node")
	// ErrSequenceMerge show that there is a sequence merge
	// indicated but got wrong values to work with
	ErrSequenceMerge = errors.New("sequence merge failed")
)

func (c config) mergeSequences(node *yaml.Node, depth uint16) error {
	// prevent loop by hardcoded limit
	if depth == c.maxDepth {
		return ErrMaxDepth
	}

	switch node.Kind {
	case yaml.DocumentNode:
		return c.mergeSequences(node.Content[0], depth+1)

	case yaml.MappingNode:
		if (len(node.Content) % 2) != 0 {
			return ErrBrokenMappingNode
		}

		for i := len(node.Content); i > 1; i = i - 2 {
			if err := c.mergeSequences(node.Content[i-1], depth+1); err != nil {
				return err
			}
		}

	case yaml.SequenceNode:
		var newContent []*yaml.Node // as long as we don't have a merge, it is empty and we don't perform slice operations
		for i := range node.Content {
			// detect "<<:" entry
			if node.Content[i].Kind == yaml.MappingNode &&
				len(node.Content[i].Content) == 2 &&
				node.Content[i].Content[0].Kind == yaml.ScalarNode &&
				node.Content[i].Content[0].Tag == mergeTag {

				// we did detect a merge tag

				if node.Content[i].Content[1].Kind == yaml.AliasNode {
					newC := node.Content[i].Content[1].Alias
					if newC.Kind != yaml.SequenceNode {
						return fmt.Errorf("%w: can only merge sequence to sequence but got something else", ErrSequenceMerge)
					}
					if len(newContent) != 0 {
						newContent = append(newContent, newC.Content...)
					} else {
						newContent = make([]*yaml.Node, i)
						copy(newContent, node.Content[:i])
						newContent = append(newContent, newC.Content...)
					}
				} else if node.Content[i].Content[1].Kind == yaml.SequenceNode {
					newC := make([]*yaml.Node, 0, len(node.Content[i].Content[1].Content))
					for _, alias := range node.Content[i].Content[1].Content {
						if alias.Kind != yaml.AliasNode {
							return fmt.Errorf("%w: merge sequences contain an non alias node", ErrSequenceMerge)
						} else if alias.Alias.Kind != yaml.SequenceNode {
							return fmt.Errorf("%w: merge sequences contain an alias to an non sequence node", ErrSequenceMerge)
						}
						newC = append(newC, alias.Alias.Content...)
					}
					if len(newContent) != 0 {
						newContent = append(newContent, newC...)
					} else {
						newContent = make([]*yaml.Node, i)
						copy(newContent, node.Content[:i])
						newContent = append(newContent, newC...)
					}
				} else {
					return fmt.Errorf("%w: merge map node did not contain an alias node", ErrSequenceMerge)
				}
			} else {

				// we did not detect a merge tag

				// else its just a normal sequence item we do process recursive
				if err := c.mergeSequences(node.Content[i], depth+1); err != nil {
					return err
				}
				// if there was a merge before we need to append it to the new content
				if len(newContent) != 0 {
					newContent = append(newContent, node.Content[i])
				}
			}
		}
		if len(newContent) != 0 {
			node.Content = newContent
		}
	}
	// we ignore Scalar and Alias Nodes
	return nil
}
