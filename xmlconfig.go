// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

//	Support load configuration from xml file
// 	Support save configuration to xml file
package config

import "encoding/xml"

// implement Coder interface
type XMLCodec struct{}

// Decode decode from byte slice which has xml format
func (xc *XMLCodec) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// Encode v to byte slice by xml Marshal
func (xc *XMLCodec) Encode(v interface{}) ([]byte, error) {
	return xml.MarshalIndent(v, "", "    ")
}

func NewXMLConfig(fileName string) *FileConfig {
	return (NewFileConfig(fileName, new(XMLCodec)))
}
