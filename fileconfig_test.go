package config

import (
	"os"
	"testing"
	"time"

	"github.com/howeyc/fsnotify"
)

const test_filename = "test.txt"

const fakeCoder_text = "test"

type fakeCoder struct {
}

func (this *fakeCoder) Encode(v interface{}) ([]byte, error) {
	return []byte(fakeCoder_text), nil
}

func (this *fakeCoder) Decode(data []byte, v interface{}) error {
	*v.(*string) = string(data)
	return nil
}

func Test_GetAndSetFileName(t *testing.T) {
	v := NewFileConfig(test_filename, new(fakeCoder))
	v.SetFileName(test_filename)
	if v.FileName() != test_filename {
		t.Errorf("FileName test failed, get %s != set %s", v.FileName(), test_filename)
	}
}

func Test_SaveAndLoadConfValue(t *testing.T) {
	fc := NewFileConfig(test_filename, new(fakeCoder))
	err := fc.SaveToFile(fakeCoder_text)
	if err != nil {
		t.Error("SavetoFile test failed", err)
	}
	defer os.Remove(test_filename)
	newvalue := fc.ConfValue()
	if newvalue != fakeCoder_text {
		t.Errorf("ConfValue test failed, get %v != set %v", newvalue, fakeCoder_text)
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
	return NA_Restart
}

func Test_FileWatcher(t *testing.T) {
	fc := NewFileConfig(test_filename, new(fakeCoder))
	err := fc.SaveToFile(fakeCoder_text)
	if err != nil {
		t.Error("SavetoFile failed ", err)
	}
	defer os.Remove(test_filename)

	err = fc.StartWatcher(watcherhandler)
	if err != nil {
		t.Error("Start watcher failed ", err)
	}
	time.Sleep(50 * time.Millisecond)
	if fc.WatcherStared() != true {
		t.Error("Start watcher failed")
	}

	var event_count = 0
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-receivedEvent:
				t.Log("Reveived event ", ev)
				event_count++
			case <-done:
				return
			case <-time.After(3 * time.Second):
				t.Error("Reveive event time out")
				return
			}
		}
	}()

	err = fc.SaveToFile(fakeCoder_text)
	if err != nil {
		t.Error("SavetoFile test failed", err)
	}

	time.Sleep(50 * time.Millisecond)
	fc.StopWatcher()
	done <- true
	if fc.WatcherStared() != false {
		t.Error("Watcher stop failed")
	}

	if event_count == 0 {
		t.Error("Receive evernt failed")
	}
}
