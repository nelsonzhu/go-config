// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

// config package is a complete configuration solution for go applications.
// It is designed to work within an application, and can handle all types of configuration needs and formats. It supports:
// 	atomic storage object
// 	reading from JSON, XML config files
// 	reading from buffer
// 	hot relaod
// 	Marshaling parser interface easy to support other format
package config

import "sync/atomic"

type Config struct {
	internalStore atomic.Value
}

func NewConfig(v interface{}) *Config {
	c := new(Config)
	c.internalStore.Store(v)
	return c
}

func (c *Config) Get() interface{} {
	return c.internalStore.Load()
}

func (c *Config) Set(v interface{}) {
	c.internalStore.Store(v)
}
