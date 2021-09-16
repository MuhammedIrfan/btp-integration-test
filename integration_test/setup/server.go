package setup

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/icon-project/btp/integration_test/setup/api"
	"github.com/labstack/echo/v4"
)

type Server struct {
	http.Server
	Port             int
	MethodRepository *api.MethodRepository
	Echo             *echo.Echo
	Ws               WSClient
}

func (s *Server) Start() {
	ln, _ := net.Listen("tcp", ":"+fmt.Sprint(s.Port))
	go func() {
		if err := s.Serve(ln); err != http.ErrServerClosed {
			log.Fatalf("Serve(): %v", err)
		}
	}()
}

func NewServer(port int, chain Chain, c interface{}) Server {
	e := echo.New()
	var group *echo.Group = e.Group("")
	// var ws = WSClient{
	// 	chain: chain,
	// }

	server := Server{
		Port:             port,
		MethodRepository: api.NewMethodRepository(),
		Echo:             e,
		Ws: WSClient{
			Chain: chain,
			Conn:  &websocket.Conn{},
		},
	}
	server.Handler = e

	grp := group.Group("/api")

	// w := group.Group("")
	grp.GET("/block", server.Ws.Block)
	grp.GET("/event", server.Ws.Event)
	grp.Use(api.JsonRpc(server.MethodRepository))
	grp.POST("", server.MethodRepository.Handle)

	return server

}

func (s *Server) RegisterMethods(handlers map[string]api.Handler) {
	for l, r := range handlers {
		s.MethodRepository.RegisterMethod(l, r)
	}

}

func (s *Server) Stop() {
	if err := s.Close(); err != nil {
		panic(err)
	}
}
