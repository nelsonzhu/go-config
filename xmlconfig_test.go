package config

import (
	"os"
	"testing"
)

const fileName_xml = "test.xml"

func Test_SaveAndLoadFromXMLFile(t *testing.T) {
	jx := NewXMLConfig(fileName_xml)
	value := test_item

	jxValue := new(testTConfig)

	err := jx.SaveToFile(value)
	if err != nil {
		t.Error("save file fail, test does NOT passed", err)
	}
	defer os.Remove(fileName_xml)

	err = jx.LoadFromFile(jxValue)
	if err != nil {
		t.Error("test does NOT passed", err)
	}
	if !(*jxValue == value && value == jx.ConfValue()) {
		t.Error("does NOT passed save != get", *jxValue, value, jx.ConfValue())
	}
}
