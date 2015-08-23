// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

// config package privide:
// 		atomic storage object
//		file configuration save, load and watcher
//		xml and Json encoder, decoder interface

package config

import "encoding/json"

// JsonConfig support load configuration from  Json file
// JsonConfig support save configuration to Json file

// implement Coder interface
type JsonCoder struct{}

// Decode decode from byte slice which has Json format
func (this *JsonCoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Encode v to byte slice by json Marshal
func (this *JsonCoder) Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

// NewJsonConfig create FileConfig support Json encode and decode
func NewJsonConfig(fileName string) *FileConfig {
	return NewFileConfig(fileName, new(JsonCoder))
}
