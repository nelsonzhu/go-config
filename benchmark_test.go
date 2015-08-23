package config

import (
	"os"
	"testing"
)

func BenchmarkConfValue(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			for i, data := range testConfigDates {
				conf := NewConfig(data)
				v := conf.ConfValue()
				if data != v {
					b.Errorf("get and set error, case:%d set:%#v get:%#v", i, data, v)
				}
			}
		}()
	}
}

func BenchmarkSaveLoad(b *testing.B) {
	jx := NewXMLConfig(fileName_json)
	value := test_item
	jxValue := new(testTConfig)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := jx.SaveToFile(value)
		if err != nil {
			b.Error("save file fail, test does NOT passed", err)
		}
		defer os.Remove(fileName_xml)

		err = jx.LoadFromFile(jxValue)
		if err != nil {
			b.Error("test does NOT passed", err)
		}
		if !(*jxValue == value && value == jx.ConfValue()) {
			b.Error("does NOT passed save != get", *jxValue, value, jx.ConfValue())
		}
	}
}
