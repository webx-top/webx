package webx

const VERSION = `1.0.0`

var (
	defaultServName string  = "webx"
	serv            *Server = NewServer(defaultServName, nil)
	servs           Servers = Servers{}
)

type Servers map[string]*Server

func (s Servers) Get(name string) (sv *Server) {
	sv, _ = s[name]
	return
}

func (s Servers) Set(name string, sv *Server) {
	s[name] = sv
}

func Serv(args ...string) (s *Server) {
	if len(args) > 0 {
		if sv, ok := servs[args[0]]; ok {
			s = sv
			return
		}
	}
	s = serv
	return
}
