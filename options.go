// SPDX-FileCopyrightText: 2023 6543 <6543@obermui.de>
// SPDX-License-Identifier: MIT

package xyaml

var defaultConfig = config{
	maxDepth: 10,
}

func NewParser(opts ...Option) Parser {
	c := defaultConfig
	for _, o := range opts {
		o(&c)
	}
	return c
}

type Option func(*config)

func WithDepth(depth uint16) Option {
	return func(c *config) {
		c.maxDepth = depth
	}
}
