package router

import (
	"net/http"
	"strings"

	"github.com/spf13/afero"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	m "github.com/rickbassham/goapi/middleware"
)

type router struct {
	address string
	r       *chi.Mux
	logger  m.RequestLogger
}

type RouteCreater interface {
	CreateRoutes(r chi.Router) chi.Router
}

type Router interface {
	ListenAndServe() error
	StaticFiles(path string, root afero.Fs)
}

func NewRouter(address string, logger m.RequestLogger, routeCreater RouteCreater) Router {
	r := chi.NewRouter()

	recoverer := m.NewRecoverer(logger)

	r.Group(func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"status": "ok"}`))
		})
	})

	r.Group(func(r chi.Router) {
		cors := cors.New(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		r.Use(cors.Handler)

		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.DefaultCompress)
		r.Use(logger.Handler())
		r.Use(recoverer.Handler())

		limiter := tollbooth.NewLimiter(10, nil)

		r.Use(tollbooth_chi.LimitHandler(limiter))

		r.Route("/api", func(r chi.Router) {
			routeCreater.CreateRoutes(r)
		})
	})

	return &router{
		address: address,
		r:       r,
		logger:  logger,
	}
}

func (svc *router) StaticFiles(path string, root afero.Fs) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	httpFs := afero.NewHttpFs(root)
	fs := http.StripPrefix(path, http.FileServer(httpFs.Dir("/")))

	if path != "/" && path[len(path)-1] != '/' {
		svc.r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	svc.r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (svc *router) ListenAndServe() error {
	return http.ListenAndServe(svc.address, svc.r)
}
