// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

// config package privide:
// 		atomic storage object
//		file configuration save, load and watcher
//		xml and Json encoder, decoder interface
package config

import (
	"errors"
	"io/ioutil"
	"reflect"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
)

type NotifyAction int

const (
	NA_Continue = iota
	NA_Restart
	NA_Stop
)

type NotifyHandler func(Sender *FileConfig, event *fsnotify.FileEvent) NotifyAction

var (
	implError     = errors.New("imtemplment of coder is wrong, cant return nil when not error")
	coderNilError = errors.New("config Coder is nil")
	saveNilError  = errors.New("can not load or save nil configuration value")
)

// Coder interface define encode and decode function
// Encode encodes configuration object to byte slice
// Decode decodes byte slice to configuration object
type Coder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// A FileConfig which is used for both reading and writing configration file.
// base is Config object which provide storage. Writing and reading simultaneously are safety.
// The coder is Coder interface which encode and decode configration object
type FileConfig struct {
	base          Config
	fileName      string
	coder         Coder
	lk            sync.RWMutex
	done          chan bool
	watcherStared bool
}

// NewFileConfig creates FileConfig object
func NewFileConfig(fileName string, coder Coder) *FileConfig {
	fc := new(FileConfig)
	fc.fileName = fileName
	fc.coder = coder
	fc.done = make(chan bool)
	return fc
}

// WatcherStared is a perperty which shows the file watcher is started or not
func (fc *FileConfig) WatcherStared() bool {
	fc.lk.RLock()
	defer fc.lk.RUnlock()
	return fc.watcherStared
}

// StopWatcher stop the file watcher
// the StopWatcher should be call when you no longer need fc object if WatcherStared
func (fc *FileConfig) StopWatcher() {
	fc.lk.Lock()
	defer fc.lk.Unlock()

	if fc.watcherStared == true {
		fc.done <- true
		select {
		case <-fc.done:
		case <-time.After(2 * time.Second):
		}
	}
	return
}

// StartWatcher start a file watcher to notify file changed
// StartWatcher call the handle once the file changed
func (fc *FileConfig) StartWatcher(handler NotifyHandler) error {
	fc.lk.Lock()
	defer fc.lk.Unlock()

	// access instence var directly
	if fc.watcherStared {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// access instance var directly
	err = watcher.Watch(fc.fileName)
	if err != nil {
		watcher.Close()
		return err
	}

	go func() {
		fc.watcherStared = true
		defer func() {
			fc.watcherStared = false
		}()
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Event:
				if handler != nil {
					action := handler(fc, ev)
					switch action {
					case NA_Stop:
						return
					case NA_Restart:
						watcher.RemoveWatch(fc.FileName())
						watcher.Watch(fc.FileName())
					}
				}
			case <-watcher.Error:
				return
			case <-fc.done:
				fc.done <- true
				return
			}
		}
	}()

	return nil
}

// LoadFromFile load content from file and decodes it to object
// save the decoded object to internal storage
// conf: decoded object receiveer, must be a pointer to object
// You should call fc function after file changed to refrash the data;
// otherwise you can call ConfValue to get object from internal storage
func (fc *FileConfig) LoadFromFile(conf interface{}) error {
	fc.lk.RLock()
	defer fc.lk.RUnlock()

	if conf == nil {
		return saveNilError
	}

	if fc.coder == nil {
		return coderNilError
	}

	data, err := ioutil.ReadFile(fc.fileName)
	if err != nil {
		return err
	}

	err = fc.coder.Decode(data, conf)
	if err != nil {
		return err
	}
	if conf == nil {
		return implError
	}

	fc.base.Set(reflect.ValueOf(conf).Elem().Interface())
	return nil
}

// SaveToFile encode conf object to byte slice and save encoded object to file
// encoded object save to internal storage simultaneously.
// writes encoded data to a file named by FileName property.
// If the file does not exist, SaveToFile creates it; otherwise truncates it before writing.
func (fc *FileConfig) SaveToFile(conf interface{}) error {
	fc.lk.Lock()
	defer fc.lk.Unlock()

	if conf == nil {
		return saveNilError
	}

	data, err := fc.coder.Encode(conf)
	if err != nil {
		return err
	}
	fc.base.Set(conf)

	return ioutil.WriteFile(fc.fileName, data, 0666)
}

// FileName property getter
func (fc *FileConfig) FileName() string {
	fc.lk.RLock()
	defer fc.lk.RUnlock()
	return fc.fileName
}

// SetFileName set config filename
func (fc *FileConfig) SetFileName(name string) {
	fc.lk.Lock()
	defer fc.lk.Unlock()

	if fc.fileName != name {
		fc.fileName = name
	}
}

// ConfValue get config value from internal storage
func (fc *FileConfig) ConfValue() interface{} {
	return fc.base.Get()
}
