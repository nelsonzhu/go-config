package config

import (
	"os"
	"testing"
)

const jasonFileName = "test.json"

func Test_SaveAndLoadFromJsonFile(t *testing.T) {
	jc := NewJSONConfig(jasonFileName)
	value := testItem

	jcValue := new(testTConfig)

	err := jc.SaveToFile(value)
	if err != nil {
		t.Error("Save file failed ", err)
	}
	defer os.Remove(jasonFileName)

	err = jc.LoadFromFile(jcValue)
	if err != nil {
		t.Error("LoadFromFile failed ", err)
	}
	if !(*jcValue == value && value == jc.Value()) {
		t.Errorf("Value to be save:%v Loaded value:%v Get saved value:%v", value, *jcValue, jc.Value())
	}
}
