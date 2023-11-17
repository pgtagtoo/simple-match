package api

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/genepg/simple-match-api/handlers"
	"github.com/genepg/simple-match-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

type Server struct {
	Router      *chi.Mux
	PersonStore *store.PersonMemoryStore
}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.PersonStore = store.NewPersonMemoryStore()

	return s
}

func (s *Server) MountHandlers() {
	flag.Parse()

	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	// Mount all handlers here
	var handler = handlers.NewHandler(s.PersonStore)

	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is simple match"))
	})

	s.Router.Post("/addPersonAndMatch", handler.AddSinglePersonAndMatch)
	s.Router.Delete("/removePerson/{id}", handler.RemoveSinglePerson)
	s.Router.Get("/queryPeople", handler.QuerySinglePeople)

	// -routes flag
	if *routes {
		// generate doc in json
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(s.Router, docgen.MarkdownOpts{
			ProjectPath: "github.com/genepg/simple-match-api",
			Intro:       "API Documentation",
		}))
		return
	}
}
