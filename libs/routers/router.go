package routers

import (
	"net/http"
)

// Router interface
type Router interface {
	Get(uri string, f func(w http.ResponseWriter, r *http.Request))
	Post(uri string, f func(w http.ResponseWriter, r *http.Request))
	Put(uri string, f func(w http.ResponseWriter, r *http.Request))
	Del(uri string, f func(w http.ResponseWriter, r *http.Request))
	Use(wmf ...func(handler http.Handler) http.Handler)
	Run(Port int, serviceName string)
}
