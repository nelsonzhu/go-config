// Copyright 2012 The go-config Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Author Nelson

// 	Reading from buffer
// 	Hot relaod
// 	Marshaling parser interface easy to support other format
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
	implError     = errors.New("imtemplment of codec is wrong, cant return nil when not error")
	codecNilError = errors.New("config Codec is nil")
	saveNilError  = errors.New("can not load or save nil configuration value")
)

// Codec interface define encode and decode function
// Encode encodes configuration object to byte slice
// Decode decodes byte slice to configuration object
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// A FileConfig which is used for both reading and writing configration file.
// base is Config object which provide storage. Writing and reading simultaneously are safety.
// The codec is Codec interface which encode and decode configration object
type FileConfig struct {
	base          Config
	fileName      string
	codec         Codec
	lk            sync.RWMutex
	done          chan bool
	watcherStared bool
}

// NewFileConfig creates FileConfig object
func NewFileConfig(fileName string, codec Codec) *FileConfig {
	fc := new(FileConfig)
	fc.fileName = fileName
	fc.codec = codec
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
// v: decoded object receiveer, must be a pointer to object
// You should call fc function after file changed to refrash the data;
// otherwise you can call Value to get object from internal storage
func (fc *FileConfig) LoadFromFile(v interface{}) error {
	fc.lk.RLock()
	defer fc.lk.RUnlock()

	if v == nil {
		return saveNilError
	}

	if fc.codec == nil {
		return codecNilError
	}

	data, err := ioutil.ReadFile(fc.fileName)
	if err != nil {
		return err
	}

	err = fc.codec.Decode(data, v)
	if err != nil {
		return err
	}
	if v == nil {
		return implError
	}

	fc.base.Set(reflect.ValueOf(v).Elem().Interface())
	return nil
}

// SaveToFile encode v object to byte slice and save encoded object to file
// encoded object save to internal storage simultaneously.
// writes encoded data to a file named by FileName property.
// If the file does not exist, SaveToFile creates it; otherwise truncates it before writing.
func (fc *FileConfig) SaveToFile(v interface{}) error {
	fc.lk.Lock()
	defer fc.lk.Unlock()

	if v == nil {
		return saveNilError
	}

	data, err := fc.codec.Encode(v)
	if err != nil {
		return err
	}
	fc.base.Set(v)

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

// Value get config value from internal storage
func (fc *FileConfig) Value() interface{} {
	return fc.base.Get()
}
