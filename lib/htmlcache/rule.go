package htmlcache

import (
	"fmt"
	"github.com/admpub/echo"
	"net/http"
	"time"
)

type Rule struct {
	SaveFile   string                                               //保存名称
	SaveFunc   func(saveFile string, c *echo.Context) string        //自定义保存名称
	ExpireTime int                                                  //过期时间(秒)
	ExpireFunc func(saveFile string, c *echo.Context) (int64, bool) //判断缓存是否过期
}

func HttpCache(ctx *echo.Context, eTag interface{}, etagValidator func(oldEtag, newEtag string) bool) bool {
	var etag string
	if eTag == nil {
		etag = fmt.Sprintf(`%v`, time.Now().UTC().Unix())
	} else {
		etag = fmt.Sprintf(`%v`, eTag)
	}
	resp := ctx.Response()
	//resp.Header().Set(`Connection`, `keep-alive`)
	resp.Header().Set(`X-Cache`, `HIT from WebX-Page-Cache`)
	if inm := ctx.Request().Header.Get("If-None-Match"); inm != `` {
		var valid bool
		if etagValidator != nil {
			valid = etagValidator(inm, etag)
		} else {
			valid = inm == etag
		}
		if valid {
			resp.Header().Del(`Content-Type`)
			resp.Header().Del(`Content-Length`)
			resp.WriteHeader(http.StatusNotModified)
			ctx.Echo().Logger().Debug(`%v is not modified.`, ctx.Path())
			return true
		}
	}
	resp.Header().Set(`Etag`, etag)
	resp.Header().Set(`Cache-Control`, `public,max-age=1`)
	return false
}
