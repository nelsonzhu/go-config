# Configuration file for Go

Cross platform: Windows, Linux, BSD and OS X.

## What is go-config?

go-config is a complete configuration solution for go applications. It is designed to work within an application, and can handle all types of configuration needs and formats. It supports:

* reading from JSON, XML config files
* reading from buffer
* hot relaod
* Marshaling interface easy support other format

## Install
			
		go get github.com/nelsonzhu/go-config

## Example:

* config

	 config  implements configuration get and set, it safe for concureency

	 Basic examples:
	
		  type myConfg struct {
			  Name     string
			  Age      int
			  Birthday time.Time
		  }
	
		  //new configuaretion service	 
		  c := NewConfig(myConfig{"test name",20,time.Now()})
	
		  //get configuration value
		  v := c.ConfValue()
	
		  //set configuration value
		  c.ConfValue(myConfig{"test name",20,time.Now()})

* fileconfig

	* provide file save, load and hot configuration reload
	* watcher function by using fsnotify package. [github.com/go-fsnotify/fsnotify](https://github.com/go-fsnotify/fsnotify)	
	Copyright (c) 2012 The Go Authors. All rights reserved.
	Copyright (c) 2012 fsnotify Authors. All rights reserved.


* xml and Json configfile
    
    provide xml and Json decoder and encoder

    base example

		var AppConfig *config.FileConfig

		type ConfigDate struct {
			Addr         string
			DBDSN        string `xml:"DB>DSN"`
			ConnPool     int    `xml:"DB>ConnPoolSize"`
			EnableSQLLog bool   `xml:"DB>EnableSQLLog"`
			SQLLogFile   string `xml:"DB>SQLLogFile"`
		}

		var test_data ConfigDate = ConfigDate{
			Addr:         "test address",
			DBDSN:        "localhost:3306",
			ConnPool:     30,
			EnableSQLLog: true,
			SQLLogFile:   "sql.log",
		}

		func configNotifyHandler(Sender *config.FileConfig, event *fsnotify.FileEvent) config.NotifyAction {
			reloadConfig()
			return config.NA_Restart
		}

		func reloadConfig() error {
			return AppConfig.LoadFromFile(new(ConfigDate))
		}

		func main() {
			var path string = "app.conf"

			AppConfig = config.NewXMLConfig(path)
			err := AppConfig.SaveToFile(test_data)
			if err != nil {
				fmt.Println("reloadConfig failed ", err)
				os.Exit(1)
			}

			err = reloadConfig()
			if err != nil {
				fmt.Println("reloadConfig failed ", err)
				os.Exit(1)
			}

			err = AppConfig.StartWatcher(configNotifyHandler)
			if err != nil {
				fmt.Println("StartWatcher error ", err)
				os.Exit(1)
			}

			fmt.Println("app conf ", AppConfig.ConfValue())

			AppConfig.StopWatcher()
		}

## License

Copyright 2012 The go-config Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.


