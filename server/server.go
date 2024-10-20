package server

import (
	"fmt"
	"net/http"
	"test/handlers"
)

type Server struct {
	HttpServer *http.Server
	Handler    handlers.HandleProcesser
}

func NewSrv() Server {
	h := handlers.NewHandler()
	srv := new(http.Server)
	return Server{HttpServer: srv, Handler: h}
}

func (s Server) Run(port string) {

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/nextdate", s.Handler.HandleDate)

	http.HandleFunc("/api/task", s.Handler.HandleTask)

	fmt.Println("Server starting at", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
