package webx

var (
	defaultServName         = "webx"
	serv            *Server = NewServer(defaultServName)
	servs           Servers = Servers{}
)

type Servers map[string]*Server

func (s Servers) Init(name string) {
	s = make(map[string]*Server)
}

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
