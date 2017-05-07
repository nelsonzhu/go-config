package config

import (
	"os"
	"testing"
)

const xmlFileName = "test.xml"

func Test_SaveAndLoadFromXMLFile(t *testing.T) {
	jx := NewXMLConfig(xmlFileName)
	value := testItem

	jxValue := new(testTConfig)

	err := jx.SaveToFile(value)
	if err != nil {
		t.Error("Save file failed ", err)
	}
	defer os.Remove(xmlFileName)

	err = jx.LoadFromFile(jxValue)
	if err != nil {
		t.Error("LoadFromFile failed ", err)
	}
	if !(*jxValue == value && value == jx.Value()) {
		t.Errorf("Value to be save:%v Loaded value:%v Get saved value:%v", value, *jxValue, jx.Value())
	}
}
