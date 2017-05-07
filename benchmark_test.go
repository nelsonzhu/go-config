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
				v := conf.Get()
				if data != v {
					b.Errorf("get and set error, case:%d set:%#v get:%#v", i, data, v)
				}
			}
		}()
	}
}

func BenchmarkSaveLoad(b *testing.B) {
	jx := NewXMLConfig(jasonFileName)
	value := testItem
	jxValue := new(testTConfig)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := jx.SaveToFile(value)
		if err != nil {
			b.Error("save file failed", err)
		}
		defer os.Remove(xmlFileName)

		err = jx.LoadFromFile(jxValue)
		if err != nil {
			b.Error("LoadFromFile failed", err)
		}
		if !(*jxValue == value && value == jx.Value()) {
			b.Error("saved values != geted", *jxValue, value, jx.Value())
		}
	}
}
