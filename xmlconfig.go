// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

package config

import "encoding/xml"

// XMLConfig support load configuration from xml file
// XMLConfig support save configuration to xml file

// implement Coder interface
type XMLCoder struct{}

// Decode decode from byte slice which has xml format
func (xc *XMLCoder) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// Encode v to byte slice by xml Marshal
func (xc *XMLCoder) Encode(v interface{}) ([]byte, error) {
	return xml.MarshalIndent(v, "", "    ")
}

func NewXMLConfig(fileName string) *FileConfig {
	return (NewFileConfig(fileName, new(XMLCoder)))
}
