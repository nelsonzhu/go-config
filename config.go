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
	confStore atomic.Value
}

func NewConfig(conf interface{}) *Config {
	c := new(Config)
	c.SetConfValue(conf)
	return c
}

func (this *Config) ConfValue() interface{} {
	return this.confStore.Load()
}

func (this *Config) SetConfValue(conf interface{}) {
	this.confStore.Store(conf)
}
