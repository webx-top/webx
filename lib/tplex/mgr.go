package tplex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/labstack/gommon/log"
)

type TemplateMgr struct {
	Caches           map[string][]byte
	Mutex            *sync.Mutex
	RootDir          string
	NewRoorDir       string
	Ignores          map[string]bool
	CachedAllows     map[string]bool
	IsReload         bool
	Logger           *log.Logger
	Preprocessor     func([]byte) []byte
	timerCallback    func() bool
	TimerCallback    func() bool
	initialized      bool
	OnChangeCallback func(string, string, string) //参数为：目标名称，类型(file/dir)，事件名(create/delete/modify/rename)
	done             chan bool
}

func (self *TemplateMgr) CloseMoniter() {
	close(self.done)
}

func (self *TemplateMgr) AllowCached(name string) bool {
	_, ok := self.CachedAllows["*.*"]
	if !ok {
		_, ok = self.CachedAllows[`*`+filepath.Ext(name)]
		if !ok {
			ok = self.CachedAllows[filepath.Base(name)]
		}
	}
	return ok
}

func (self *TemplateMgr) Moniter(rootDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	//fmt.Println("[webx] TemplateMgr watcher is start.")
	defer watcher.Close()
	self.done = make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					break
				}
				if _, ok := self.Ignores[filepath.Base(ev.Name)]; ok {
					break
				}
				if _, ok := self.Ignores[`*`+filepath.Ext(ev.Name)]; ok {
					break
				}
				d, err := os.Stat(ev.Name)
				if err != nil {
					break
				}

				if ev.IsCreate() {
					if d.IsDir() {
						watcher.Watch(ev.Name)
						self.OnChange(ev.Name, "dir", "create")
					} else {
						self.OnChange(ev.Name, "file", "create")
						if self.AllowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							content, err := ioutil.ReadFile(ev.Name)
							if err != nil {
								self.Logger.Info("loaded template %v failed: %v", tmpl, err)
								break
							}
							self.Logger.Info("loaded template file %v success", tmpl)
							self.CacheTemplate(tmpl, content)
						}
					}
				} else if ev.IsDelete() {
					if d.IsDir() {
						watcher.RemoveWatch(ev.Name)
						self.OnChange(ev.Name, "dir", "delete")
					} else {
						self.OnChange(ev.Name, "file", "delete")
						if self.AllowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							self.CacheDelete(tmpl)
						}
					}
				} else if ev.IsModify() {
					if d.IsDir() {
						self.OnChange(ev.Name, "dir", "modify")
					} else {
						self.OnChange(ev.Name, "file", "modify")
						if self.AllowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							content, err := ioutil.ReadFile(ev.Name)
							if err != nil {
								self.Logger.Error("reloaded template %v failed: %v", tmpl, err)
								break
							}
							self.CacheTemplate(tmpl, content)
							self.Logger.Info("reloaded template %v success", tmpl)
						}
					}
				} else if ev.IsRename() {
					if d.IsDir() {
						watcher.RemoveWatch(ev.Name)
						self.OnChange(ev.Name, "dir", "rename")
					} else {
						self.OnChange(ev.Name, "file", "rename")
						if self.AllowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							self.CacheDelete(tmpl)
						}
					}
				}
			case err := <-watcher.Error:
				self.Logger.Error("error:", err)
			case <-time.After(time.Second * 2):
				if self.timerCallback != nil {
					if self.timerCallback() == false {
						close(self.done)
						return
					}
				}
				//fmt.Printf("TemplateMgr timer operation: %v.\n", time.Now())
			}
		}
	}()

	err = filepath.Walk(self.RootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Watch(f)
		}
		return nil
	})

	if err != nil {
		self.Logger.Error(err.Error())
		return err
	}

	<-self.done
	//fmt.Println("[webx] TemplateMgr watcher is closed.")
	return nil
}

func (self *TemplateMgr) OnChange(name, typ, event string) {
	if self.OnChangeCallback != nil {
		self.Mutex.Lock()
		defer self.Mutex.Unlock()
		name = FixDirSeparator(name)
		self.OnChangeCallback(name[len(self.RootDir)+1:], typ, event)
	}
}

func (self *TemplateMgr) CacheAll(rootDir string) error {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	fmt.Print("Reading the contents of the template files, please wait... ")
	err := filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tmpl := f[len(rootDir)+1:]
		tmpl = FixDirSeparator(tmpl)
		if _, ok := self.Ignores[filepath.Base(tmpl)]; !ok {
			fpath := filepath.Join(self.RootDir, tmpl)
			content, err := ioutil.ReadFile(fpath)
			if err != nil {
				self.Logger.Debug("load template %s error: %v", fpath, err)
				return err
			}
			self.Logger.Debug("loaded template", fpath)
			self.Caches[tmpl] = content
		}
		return nil
	})
	fmt.Println("Complete.")
	return err
}

func (self *TemplateMgr) defaultTimerCallback() func() bool {
	return func() bool {
		if self.TimerCallback != nil {
			return self.TimerCallback()
		}
		//更改模板主题后，关闭当前监控，重新监控新目录
		if self.NewRoorDir == "" || self.NewRoorDir == self.RootDir {
			return true
		}
		self.ClearCache()
		self.Ignores = make(map[string]bool)
		self.RootDir = self.NewRoorDir
		go self.Moniter(self.RootDir)
		return false
	}
}

func (self *TemplateMgr) Close() {
	self.TimerCallback = func() bool {
		self.ClearCache()
		self.Ignores = make(map[string]bool)
		self.TimerCallback = nil
		return false
	}
	self.initialized = false
}

func (self *TemplateMgr) Init(logger *log.Logger, rootDir string, reload bool, allows ...string) error {
	if self.initialized {
		if rootDir == self.RootDir {
			return nil
		} else {
			self.TimerCallback = func() bool {
				self.ClearCache()
				self.Ignores = make(map[string]bool)
				self.CachedAllows = make(map[string]bool)
				self.TimerCallback = nil
				return false
			}
		}
	} else if !reload {
		self.TimerCallback = func() bool {
			self.TimerCallback = nil
			return false
		}
	}
	self.RootDir = rootDir
	self.Caches = make(map[string][]byte)
	self.Ignores = make(map[string]bool)
	self.CachedAllows = make(map[string]bool)
	for _, allow := range allows {
		self.CachedAllows[allow] = true
	}
	self.Mutex = &sync.Mutex{}
	self.Logger = logger
	if dirExists(rootDir) {
		//self.CacheAll(rootDir)
		if reload {
			self.timerCallback = self.defaultTimerCallback()
			go self.Moniter(rootDir)
		}
	}

	if len(self.Ignores) == 0 {
		self.Ignores["*.tmp"] = false
		self.Ignores["*.TMP"] = false
	}
	if len(self.CachedAllows) == 0 {
		self.CachedAllows["*.*"] = true
	}
	self.initialized = true
	return nil
}

func (self *TemplateMgr) GetTemplate(tmpl string) ([]byte, error) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	tmpl = FixDirSeparator(tmpl)
	if tmpl[0] == '/' {
		tmpl = tmpl[1:]
	}

	if content, ok := self.Caches[tmpl]; ok {
		self.Logger.Debug("load template %v from cache", tmpl)
		return content, nil
	}

	content, err := ioutil.ReadFile(filepath.Join(self.RootDir, tmpl))
	if err == nil {
		self.Logger.Debug("load template %v from the file:", tmpl)
		self.Caches[tmpl] = content
	}
	return content, err
}

func (self *TemplateMgr) CacheTemplate(tmpl string, content []byte) {
	if self.Preprocessor != nil {
		content = self.Preprocessor(content)
	}
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	tmpl = FixDirSeparator(tmpl)
	self.Logger.Debug("update template %v on cache", tmpl)
	self.Caches[tmpl] = content
	return
}

func (self *TemplateMgr) CacheDelete(tmpl string) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()
	tmpl = FixDirSeparator(tmpl)
	if _, ok := self.Caches[tmpl]; ok {
		self.Logger.Info("delete template %v from cache", tmpl)
		delete(self.Caches, tmpl)
	}
	return
}

func (self *TemplateMgr) ClearCache() {
	self.Caches = make(map[string][]byte)
}
