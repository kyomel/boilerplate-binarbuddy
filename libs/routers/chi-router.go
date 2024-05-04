package routers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chiRouter struct{}

var (
	chiDispatcher = chi.NewRouter()
)

// NewChiRouter func
func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) Get(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

func (*chiRouter) Post(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

func (*chiRouter) Put(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Put(uri, f)
}

func (*chiRouter) Del(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Delete(uri, f)
}

func (*chiRouter) Use(mwf ...func(handler http.Handler) http.Handler) {
	for _, h := range mwf {
		chiDispatcher.Use(h)
	}
}

func (*chiRouter) Run(port int, serviceName string) {
	log.Println(fmt.Sprintf(serviceName+" - Chi HTTP Server was running on port %d...\n", port))
	go http.ListenAndServe(fmt.Sprintf(":%d", port), chiDispatcher)
	select {}
}
