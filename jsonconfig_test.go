package config

import (
	"os"
	"testing"
)

const fileName_json = "test.json"

func Test_SaveAndLoadFromJsonFile(t *testing.T) {
	jc := NewJsonConfig(fileName_json)
	value := test_item

	jcValue := new(testTConfig)

	err := jc.SaveToFile(value)
	if err != nil {
		t.Error("Save file failed ", err)
	}
	defer os.Remove(fileName_json)

	err = jc.LoadFromFile(jcValue)
	if err != nil {
		t.Error("LoadFromFile failed ", err)
	}
	if !(*jcValue == value && value == jc.ConfValue()) {
		t.Error("Saved value != geted ", *jcValue, value, jc.ConfValue())
	}
}
