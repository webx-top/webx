package com

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/howeyc/fsnotify"
)

//监控时间函数
type MoniterEventFunc struct {
	Create func(string) //创建
	Delete func(string) //删除
	Modify func(string) //修改
	Rename func(string) //重命名
	Timer  func() bool  //定时操作
}

//文件监测
func Moniter(rootDir string, callback MoniterEventFunc, filter func(string) bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					break
				}
				if filter != nil {
					if !filter(ev.Name) {
						break
					}
				}
				d, err := os.Stat(ev.Name)
				if err != nil {
					break
				}

				if callback.Create != nil && ev.IsCreate() {
					if d.IsDir() {
						watcher.Watch(ev.Name)
					} else {
						callback.Create(ev.Name)
					}
				} else if callback.Delete != nil && ev.IsDelete() {
					if d.IsDir() {
						watcher.RemoveWatch(ev.Name)
					} else {
						callback.Delete(ev.Name)
					}
				} else if callback.Modify != nil && ev.IsModify() {
					if d.IsDir() {
					} else {
						callback.Modify(ev.Name)
					}
				} else if callback.Rename != nil && ev.IsRename() {
					if d.IsDir() {
						watcher.RemoveWatch(ev.Name)
					} else {
						callback.Rename(ev.Name)
					}
				}
			case err := <-watcher.Error:
				fmt.Println("Watcher error:", err)
			case <-time.After(time.Second * 2):
				if callback.Timer != nil {
					if callback.Timer() == false {
						close(done)
						return
					}
				}
				//fmt.Printf("Moniter timer operation: %v.\n", time.Now())
			}
		}
	}()

	err = filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Watch(f)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	<-done
	return nil
}
