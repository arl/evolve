// +build dev

package main

import (
	"log"

	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

func init() {
	go liveReload()
}

func liveReload() {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					lr.Reload(event.Name)
				}
			case err := <-watcher.Errors:
				lr.Alert(err.Error())
			}
		}
	}()

	err = watcher.Add("index.html")
	if err != nil {
		log.Fatalln("liveReload watcher:", err)
	}

	select {}
}
