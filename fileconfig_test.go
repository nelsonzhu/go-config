package config

import (
	"os"
	"testing"
	"time"

	"github.com/howeyc/fsnotify"
)

const testFilename = "test.txt"

const fakeCoderText = "test"

type fakeCoder struct {
}

func (fc *fakeCoder) Encode(v interface{}) ([]byte, error) {
	return []byte(fakeCoderText), nil
}

func (fc *fakeCoder) Decode(data []byte, v interface{}) error {
	*v.(*string) = string(data)
	return nil
}

func Test_GetAndSetFileName(t *testing.T) {
	v := NewFileConfig(testFilename, new(fakeCoder))
	v.SetFileName(testFilename)
	if v.FileName() != testFilename {
		t.Errorf("FileName test failed, get %s != set %s", v.FileName(), testFilename)
	}
}

func Test_SaveAndLoadConfValue(t *testing.T) {
	fc := NewFileConfig(testFilename, new(fakeCoder))
	err := fc.SaveToFile(fakeCoderText)
	if err != nil {
		t.Error("SavetoFile test failed", err)
	}
	defer os.Remove(testFilename)
	newvalue := fc.Value()
	if newvalue != fakeCoderText {
		t.Errorf("ConfValue test failed, get %v != set %v", newvalue, fakeCoderText)
	}
	str := new(string)
	err = fc.LoadFromFile(str)
	if err != nil {
		t.Errorf("LoadFromFile failed %v", err)
	}
	if *str != newvalue {
		t.Errorf("LoadFromFile test failed, get %v != set %v", *str, newvalue)
	}
}

var receivedEvent = make(chan *fsnotify.FileEvent)

func watcherhandler(Sender *FileConfig, event *fsnotify.FileEvent) NotifyAction {
	receivedEvent <- event
	return NARestart
}

func Test_FileWatcher(t *testing.T) {
	fc := NewFileConfig(testFilename, new(fakeCoder))
	err := fc.SaveToFile(fakeCoderText)
	if err != nil {
		t.Error("SavetoFile failed ", err)
	}
	defer os.Remove(testFilename)

	err = fc.StartWatcher(watcherhandler)
	if err != nil {
		t.Error("Start watcher failed ", err)
	}
	time.Sleep(50 * time.Millisecond)
	if fc.WatcherStared() != true {
		t.Error("Start watcher failed")
	}

	var eventCount = 0
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-receivedEvent:
				t.Log("Reveived event ", ev)
				eventCount++
			case <-done:
				return
			case <-time.After(3 * time.Second):
				t.Error("Reveive event time out")
				return
			}
		}
	}()

	err = fc.SaveToFile(fakeCoderText)
	if err != nil {
		t.Error("SavetoFile test failed", err)
	}

	time.Sleep(50 * time.Millisecond)
	fc.StopWatcher()
	done <- true
	if fc.WatcherStared() != false {
		t.Error("Watcher stop failed")
	}

	if eventCount == 0 {
		t.Error("Receive evernt failed")
	}
}
