package ssi

var DefaultName = `SESSIONID`

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}
