package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct {
	Handler http.HandlerFunc
	Method  string
	Path    string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]Handler
	WebServerPort string
}

func NewWebServer(webServerPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]Handler),
		WebServerPort: webServerPort,
	}
}

func (ws *WebServer) AddHandler(path string, method string, handler http.HandlerFunc) {
	ws.Handlers[path+"-"+method] = Handler{
		Handler: handler,
		Method:  method,
		Path:    path,
	}
}

func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Logger)

	for _, handler := range ws.Handlers {
		switch handler.Method {
		case "GET":
			ws.Router.Get(handler.Path, handler.Handler)
		case "POST":
			ws.Router.Post(handler.Path, handler.Handler)
		case "PUT":
			ws.Router.Put(handler.Path, handler.Handler)
		case "DELETE":
			ws.Router.Delete(handler.Path, handler.Handler)
		default:
			ws.Router.Head(handler.Path, handler.Handler)
		}
	}

	fmt.Println("Starting web server on port", ws.WebServerPort)
	err := http.ListenAndServe(":"+ws.WebServerPort, ws.Router)
	if err != nil {
		panic(err)
	}
}
