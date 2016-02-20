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
	"fmt"

	"github.com/coscms/xorm"
	"github.com/webx-top/webx/lib/client"
	"github.com/webx-top/webx/lib/database"
)

func NewSelect(orm *database.Orm, c client.Client) *Select {
	s := &Select{
		Orm:    orm,
		Params: make([]interface{}, 0),
		Client: c,
	}
	return s
}

type Select struct {
	Offset    int64
	Limit     int64
	OrderBy   string
	Condition string
	Params    []interface{}
	GroupBy   string
	Having    string
	Table     interface{}
	Alias     string
	*database.Orm
	client.Client
}

func (a *Select) Do(args ...interface{}) *xorm.Session {
	return a.GenSess(args...).OrderBy(a.OrderBy).Limit(int(a.Limit), int(a.Offset))
}

func (a *Select) AddParam(args ...interface{}) *Select {
	a.Params = append(a.Params, args...)
	return a
}

func (a *Select) FromClient(gen bool, fields ...string) *Select {
	a.OrderBy = a.Client.GenOrderBy()
	a.Offset = a.Client.Offset()
	a.Limit = a.Client.PageSize()
	if !gen {
		return a
	}
	sql := a.Condition
	sch := a.Client.GenSearch(fields...)
	if sch != `` {
		if sql != `` {
			sql += ` AND `
		}
		sql += sch
	}
	a.Condition = sql
	return a
}

func (a *Select) GenSess(args ...interface{}) *xorm.Session {
	var s *xorm.Session = a.Orm.NewSession()
	s.IsAutoClose = true
	switch len(args) {
	case 2:
		alias, _ := args[1].(string)
		if args[0] == nil {
			s = s.Alias(alias)
		} else {
			s = s.Table(args[0]).Alias(alias)
			a.Table = args[0]
		}
		a.Alias = alias
	case 1:
		s = s.Table(args[0])
		a.Table = args[0]
	default:
		if a.Table != nil {
			s = s.Table(a.Table)
		}
		if a.Alias != `` {
			s = s.Alias(a.Alias)
		}
	}
	s = s.Where(a.Condition, a.Params...).GroupBy(a.GroupBy)
	if a.Having != `` {
		s = s.Having(a.Having)
	}
	return s
}

func (a *Select) Count(m interface{}) int64 {
	count, err := a.GenSess().Count(m)
	if err != nil {
		fmt.Println(err)
	}
	return count
}
