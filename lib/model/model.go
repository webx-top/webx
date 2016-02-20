/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package model

import (
	"github.com/coscms/xorm"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/client"
	"github.com/webx-top/webx/lib/database"
	"github.com/webx-top/webx/lib/i18n"
)

func NewModel(db *database.Orm, ctx *X.Context) *Model {
	return &Model{
		DB:      db,
		Context: ctx,
	}
}

type Model struct {
	DB      *database.Orm
	Context *X.Context
}

func (this *Model) T(key string, args ...interface{}) string {
	return i18n.T(this.Context.Language, key, args...)
}

// =====================================
// TransManager
// =====================================
func (this *Model) Begin() *xorm.Session {
	ss, _ := this.transSession()
	if ss != nil {
		ss.Close()
	}
	ss = this.DB.NewSession()
	err := ss.Begin()
	if err != nil {
		this.Context.Object().Echo().Logger().Error(err)
	}
	this.Context.Set(`webx:transSession`, ss)
	return ss
}

//事务是否已经开始
func (this *Model) HasBegun() bool {
	ss, ok := this.transSession()
	return ok && ss != nil
}

func (this *Model) transSession() (ss *xorm.Session, ok bool) {
	ss, ok = this.Context.Get(`webx:transSession`).(*xorm.Session)
	return
}

func (this *Model) TSess() *xorm.Session { // TransSession
	ss, ok := this.transSession()
	if ok == false || ss == nil {
		return this.Begin()
	}
	return ss
}

func (this *Model) Trans(fn func() error) *database.Orm {
	ss, ok := this.transSession()
	begun := ok && ss != nil
	if !begun {
		ss = this.Begin()
	}
	result := fn()
	if !begun {
		this.End(result == nil, ss)
	}
	return this.DB
}

func (this *Model) Sess() *xorm.Session { // TransSession or Session
	ss, ok := this.transSession()
	if ok == false {
		var ss *xorm.Session = this.DB.NewSession()
		ss.IsAutoClose = true
		return ss
	}
	return ss
}

func (this *Model) End(result bool, args ...*xorm.Session) (err error) {
	var ss *xorm.Session
	if len(args) > 0 && args[0] != nil {
		ss = args[0]
	} else {
		ss, _ = this.transSession()
	}
	if result {
		err = ss.Commit()
	} else {
		err = ss.Rollback()
	}
	if err != nil {
		this.Context.Object().Echo().Logger().Error(err)
	}
	ss.Close()
	this.Context.Set(`webx:transSession`, nil)
	return
}

func (this *Model) NewSelect(m interface{}) *Select {
	return NewSelect(this.DB, this.NewClient(m))
}

func (this *Model) NewClient(m interface{}) client.Client {
	clientName := this.Context.Query(`client`)
	c := client.Get(clientName)
	return c.Init(this.Context, this.DB, m)
}
