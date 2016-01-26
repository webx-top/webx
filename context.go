package webx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/cookie"
	sessLib "github.com/webx-top/webx/lib/session"
)

func NewContext(s *Server, c echo.Context) *Context {
	return &Context{
		Context: c,
		Server:  s,
	}
}

const (
	NO_PERM = -2 //无权限
	NO_LOGN = -1 //未登录
	FAILURE = 0  //操作失败
	SUCCESS = 1  //操作成功
)

type Output struct {
	Status  int
	Message interface{}
	Data    interface{}
}

type Context struct {
	echo.Context
	*Server
	*App
	*Output
	session        sessLib.Session
	C              interface{}
	ControllerName string
	ActionName     string
	Language       string
	Code           int
	Tmpl           string
	Format         string
	Exit           bool
}

func (c *Context) Init(app *App, ctl interface{}, ctlName string, actName string) {
	c.App = app
	c.C = ctl
	c.ControllerName = ctlName
	c.ActionName = actName
}

func (c *Context) Reset(r *http.Request, w http.ResponseWriter, e *echo.Echo) {
	c.Context.Reset(r, w, e)
	c.ControllerName = ``
	c.Server = nil
	c.App = nil
	c.ActionName = ``
	c.Language = ``
	c.Exit = false
	c.Output = &Output{1, ``, make(map[string]string)}
	c.Tmpl = ``
	c.Format = c.ResolveFormat()
}

func (c *Context) InitSession(session sessLib.Session) {
	if session == nil {
		session = sessLib.NewSession(c.Server.SessionStoreEngine,
			c.Server.SessionStoreConfig,
			c.Request(), c.Response())
	}
	c.session = session
}

func (c *Context) Session() sessLib.Session {
	if c.session == nil {
		c.InitSession(nil)
	}
	return c.session
}

func (c *Context) SetSession(key string, val interface{}) {
	s := c.Session()
	s.Set(key, val)
	s.Save()
}

func (c *Context) GetSession(key string) interface{} {
	return c.Session().Get(key)
}

func (c *Context) Cookie(key string, value string) *cookie.Cookie {
	liftTime := c.Server.CookieExpires
	sPath := "/"
	domain := c.Server.CookieDomain
	secure := c.IsSecure()
	httpOnly := c.Server.CookieHttpOnly
	return cookie.New(c.Server.CookiePrefix+key, value, liftTime, sPath, domain, secure, httpOnly)
}

func (c *Context) GetCookie(key string) string {
	var val string
	if res, err := c.Request().Cookie(c.Server.CookiePrefix + key); err == nil && res.Value != "" {
		val, _ = com.UrlDecode(res.Value)
	}
	return val
}

func (c *Context) SetCookie(key, val string, args ...interface{}) {
	val = com.UrlEncode(val)
	cookie := c.Cookie(key, val)
	switch len(args) {
	case 5:
		httpOnly, _ := args[4].(bool)
		cookie.HttpOnly(httpOnly)
		fallthrough
	case 4:
		secure, _ := args[3].(bool)
		cookie.Secure(secure)
		fallthrough
	case 3:
		domain, _ := args[2].(string)
		cookie.Domain(domain)
		fallthrough
	case 2:
		path, _ := args[1].(string)
		cookie.Path(path)
		fallthrough
	case 1:
		var liftTime int64
		switch args[0].(type) {
		case int:
			liftTime = int64(args[0].(int))
		case int64:
			liftTime = args[0].(int64)
		case time.Duration:
			liftTime = int64(args[0].(time.Duration))
		}
		cookie.Expires(liftTime)
	}
	cookie.Send(c)
}

func (c *Context) SetSecCookie(key string, value interface{}) {
	if c.Server.Codec == nil {
		val, _ := value.(string)
		c.SetCookie(key, val)
		return
	}
	encoded, err := c.Server.Codec.Encode(key, value)
	if err != nil {
		c.X().Echo().Logger().Error(err)
	} else {
		c.SetCookie(key, encoded)
	}
}

func (c *Context) GetSecCookie(key string) (value interface{}) {
	cookieValue := c.GetCookie(key)
	if cookieValue != "" && c.Server.Codec != nil {
		err := c.Server.Codec.Decode(key, cookieValue, &value)
		if err != nil {
			c.X().Echo().Logger().Error(err)
		}
	}
	return
}

func (c *Context) Body() ([]byte, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}
	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func (c *Context) IP() string {
	proxy := c.Proxy()
	if len(proxy) > 0 && proxy[0] != "" {
		return proxy[0]
	}
	ip := strings.Split(c.Request().RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (c *Context) Header(name string) string {
	return c.Request().Header.Get(name)
}

func (c *Context) Method() string {
	return c.Request().Method
}

func (c *Context) IsAjax() bool {
	return c.Header("X-Requested-With") == "XMLHttpRequest"
}

func (c *Context) IsPjax() bool {
	return c.Header("X-PJAX") == "true"
}

func (c *Context) PjaxContainer() string {
	return c.Header("X-PJAX-Container")
}

func (c *Context) IsOnlyAjax() bool {
	return c.IsAjax() && c.Header("X-PJAX") == ""
}

//CREATE：在服务器新建一个资源
func (c *Context) IsPost() bool {
	return c.Method() == "POST"
}

//SELECT：从服务器取出资源（一项或多项）
func (c *Context) IsGet() bool {
	return c.Method() == "GET"
}

//UPDATE：在服务器更新资源（客户端提供改变后的完整资源）
func (c *Context) IsPut() bool {
	return c.Method() == "PUT"
}

//DELETE：从服务器删除资源
func (c *Context) IsDel() bool {
	return c.Method() == "DELETE"
}

//获取资源的元数据
func (c *Context) IsHead() bool {
	return c.Method() == "HEAD"
}

//UPDATE：在服务器更新资源（客户端提供改变的属性）
func (c *Context) IsPatch() bool {
	return c.Method() == "PATCH"
}

//获取信息，关于资源的哪些属性是客户端可以改变的
func (c *Context) IsOptions() bool {
	return c.Method() == "OPTIONS"
}

// Form returns form parameter by name.
func (c *Context) Form(name string) string {
	r := c.Request()
	if r.Form == nil {
		if c.IsUpload() {
			r.ParseMultipartForm(c.Server.MaxUploadSize)
			if len(r.PostForm) == 0 {
				r.PostForm = r.MultipartForm.Value
			}
		} else {
			r.ParseForm()
		}
	}
	return c.Context.Form(name)
}

func (c *Context) IsSecure() bool {
	return c.Scheme() == "https"
}

// IsWebsocket returns boolean of this request is in webSocket.
func (c *Context) IsWebsocket() bool {
	return c.Header("Upgrade") == "websocket"
}

// IsUpload returns boolean of whether file uploads in this request or not..
func (c *Context) IsUpload() bool {
	return strings.Contains(c.Header("Content-Type"), "multipart/form-data")
}

// Get the content type.
// e.g. From "multipart/form-data; boundary=--" to "multipart/form-data"
// If none is specified, returns "text/html" by default.
func (c *Context) ResolveContentType() string {
	contentType := c.Header("Content-Type")
	if contentType == "" {
		return "text/html"
	}
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

// ResolveFormat maps the request's Accept MIME type declaration to
// a Request.Format attribute, specifically "html", "xml", "json", or "txt",
// returning a default of "html" when Accept header cannot be mapped to a
// value above.
func (c *Context) ResolveFormat() string {
	format := c.Query("format")
	if format != `` {
		return format
	}
	accept := c.Header("Accept")
	switch {
	case accept == "",
		strings.HasPrefix(accept, "*/*"), // */
		strings.Contains(accept, "application/xhtml"),
		strings.Contains(accept, "text/html"):
		return "html"
	case strings.Contains(accept, "application/json"),
		strings.Contains(accept, "text/javascript"),
		strings.Contains(accept, "application/javascript"):
		return "json"
	case strings.Contains(accept, "application/xml"),
		strings.Contains(accept, "text/xml"):
		return "xml"
	case strings.Contains(accept, "text/plain"):
		return "text"
	}

	return "html"
}

// Protocol returns request protocol name, such as HTTP/1.1 .
func (c *Context) Protocol() string {
	return c.Request().Proto
}

// Site returns base site url as scheme://domain type.
func (c *Context) Site() string {
	return c.Scheme() + "://" + c.Domain()
}

// Scheme returns request scheme as "http" or "https".
func (c *Context) Scheme() string {
	if c.Request().URL.Scheme != "" {
		return c.Request().URL.Scheme
	}
	if c.Request().TLS == nil {
		return "http"
	}
	return "https"
}

// Domain returns host name.
// Alias of Host method.
func (c *Context) Domain() string {
	return c.Host()
}

// Host returns host name.
// if no host info in request, return localhost.
func (c *Context) Host() string {
	if c.Request().Host != "" {
		hostParts := strings.Split(c.Request().Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return c.Request().Host
	}
	return "localhost"
}

// Proxy returns proxy client ips slice.
func (c *Context) Proxy() []string {
	if ips := c.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// Referer returns http referer header.
func (c *Context) Referer() string {
	return c.Header("Referer")
}

// Refer returns http referer header.
func (c *Context) Refer() string {
	return c.Referer()
}

// SubDomains returns sub domain string.
// if aa.bb.domain.com, returns aa.bb .
func (c *Context) SubDomains() string {
	parts := strings.Split(c.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

// Port returns request client port.
// when error or empty, return 80.
func (c *Context) Port() int {
	parts := strings.Split(c.Request().Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}

func (c *Context) Assign(key string, val interface{}) {
	data, _ := c.Output.Data.(map[string]interface{})
	if data == nil {
		data = map[string]interface{}{}
	}
	data[key] = val
	c.Output.Data = data
}

func (c *Context) AssignX(values *map[string]interface{}) {
	if values == nil {
		return
	}
	data, _ := c.Output.Data.(map[string]interface{})
	for key, val := range *values {
		data[key] = val
	}
	c.Output.Data = data
}

func (c *Context) Display(args ...interface{}) error {
	switch len(args) {
	case 2:
		if v, ok := args[0].(string); ok && v != `` {
			c.Tmpl = v
		}
		if v, ok := args[1].(int); ok && v > 0 {
			c.Code = v
		}
	case 1:
		if v, ok := args[0].(int); ok {
			c.Code = v
		} else if v, ok := args[0].(string); ok {
			c.Tmpl = v
		}
	}
	if c.Code <= 0 {
		c.Code = http.StatusOK
	}
	if ignore, _ := c.Get(`webx:ignoreRender`).(bool); ignore {
		return nil
	}

	switch c.Format {
	case `xml`:
		b, err := xml.Marshal(c.Output)
		if err != nil {
			return err
		}
		c.X().Xml(c.Code, b)
		return nil
	case `json`:
		b, err := json.Marshal(c.Output)
		if err != nil {
			return err
		}
		callback := c.Query(`callback`)
		if callback != `` {
			c.X().Jsonp(c.Code, callback, b)
		} else {
			c.X().Json(c.Code, b)
		}
		return nil
	default:
		if c.Tmpl == `` {
			return nil
		}
		c.Context.SetFunc(`Status`, func() int {
			return c.Output.Status
		})
		c.Context.SetFunc(`Message`, func() interface{} {
			return c.Output.Message
		})
		return c.Render(c.Code, c.Tmpl, c.Output.Data)
	}
}