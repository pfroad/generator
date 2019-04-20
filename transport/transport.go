package transport

import (
	"github.com/julienschmidt/httprouter"
)

func Start() {

}

type Transport interface {
	Start() error
}

type RestfulTransport struct {
	addr     string
	router   *httprouter.Router
	handlers map[string]httprouter.Handle
}

func NewRestfulTransport(addr string) *RestfulTransport {
	return &RestfulTransport{
		addr: addr,
	}
}

func (self *RestfulTransport) start() {

}
