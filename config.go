// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

// config package privide:
// 		atomic storage object
//		file configuration save, load and watcher
//		xml and Json encoder, decoder interface

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
