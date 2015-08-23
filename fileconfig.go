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
func (this *FileConfig) WatcherStared() bool {
	this.lk.RLock()
	defer this.lk.RUnlock()
	return this.watcherStared
}

// StopWatcher stop the file watcher
// the StopWatcher should be call when you no longer need this object if WatcherStared
func (this *FileConfig) StopWatcher() {
	this.lk.Lock()
	defer this.lk.Unlock()

	if this.watcherStared == true {
		this.done <- true
		select {
		case <-this.done:
		case <-time.After(2 * time.Second):
		}
	}
	return
}

// StartWatcher start a file watcher to notify file changed
// StartWatcher call the handle once the file changed
func (this *FileConfig) StartWatcher(handler NotifyHandler) error {
	this.lk.Lock()
	defer this.lk.Unlock()

	// access instence var directly
	if this.watcherStared {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// access instance var directly
	err = watcher.Watch(this.fileName)
	if err != nil {
		watcher.Close()
		return err
	}

	go func(w *fsnotify.Watcher) {
		this.watcherStared = true
		defer func() {
			this.watcherStared = false
		}()
		defer w.Close()
		for {
			select {
			case ev := <-w.Event:
				if handler != nil {
					action := handler(this, ev)
					switch action {
					case NA_Stop:
						return
					case NA_Restart:
						w.RemoveWatch(this.FileName())
						w.Watch(this.FileName())
					}
				}
			case <-w.Error:
				return
			case <-this.done:
				this.done <- true
				return
			}
		}
	}(watcher)

	return nil
}

// LoadFromFile load content from file and decodes it to object
// save the decoded object to internal storage
// conf: decoded object receiveer, must be a pointer to object
// You should call this function after file changed to refrash the data;
// otherwise you can call ConfValue to get object from internal storage
func (this *FileConfig) LoadFromFile(conf interface{}) error {
	this.lk.RLock()
	defer this.lk.RUnlock()

	if conf == nil {
		return saveNilError
	}

	if this.coder == nil {
		return coderNilError
	}

	data, err := ioutil.ReadFile(this.fileName)
	if err != nil {
		return err
	}

	err = this.coder.Decode(data, conf)
	if err != nil {
		return err
	}
	if conf == nil {
		return implError
	}

	this.base.SetConfValue(reflect.ValueOf(conf).Elem().Interface())
	return nil
}

// SaveToFile encode conf object to byte slice and save encoded object to file
// encoded object save to internal storage simultaneously.
// writes encoded data to a file named by FileName property.
// If the file does not exist, SaveToFile creates it; otherwise truncates it before writing.
func (this *FileConfig) SaveToFile(conf interface{}) error {
	this.lk.Lock()
	defer this.lk.Unlock()

	if conf == nil {
		return saveNilError
	}

	data, err := this.coder.Encode(conf)
	if err != nil {
		return err
	}
	this.base.SetConfValue(conf)

	return ioutil.WriteFile(this.fileName, data, 0666)
}

// FileName property getter
func (this *FileConfig) FileName() string {
	this.lk.RLock()
	defer this.lk.RUnlock()
	return this.fileName
}

// SetFileName set config filename
func (this *FileConfig) SetFileName(name string) {
	this.lk.Lock()
	defer this.lk.Unlock()

	if this.fileName != name {
		this.fileName = name
	}
}

// ConfValue get config value from internal storage
func (this *FileConfig) ConfValue() interface{} {
	return this.base.ConfValue()
}
