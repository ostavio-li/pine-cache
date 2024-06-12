package http

import (
	"io"
	"ljcgo/cache"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	// http.Handle("/status", s.statusHandler())
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		return
	}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

type cacheHandler struct {
	*Server
}

func (h cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := r.Method
	if m == http.MethodPut {
		b, _ := io.ReadAll(r.Body)
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}

	if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
		return
	}

	if m == http.MethodDelete {
		err := h.Del(key)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
