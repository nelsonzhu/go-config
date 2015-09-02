// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

//	Support load configuration from  Json file
// 	Support save configuration to Json file
package config

import "encoding/json"

// implement Codec interface
type JsonCodec struct{}

// Decode decode from byte slice which has Json format
func (jc *JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Encode v to byte slice by json Marshal
func (jc *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

// NewJsonConfig create FileConfig support Json encode and decode
func NewJsonConfig(fileName string) *FileConfig {
	return NewFileConfig(fileName, new(JsonCodec))
}
