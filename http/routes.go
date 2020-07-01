package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"net/http"
)

func RegisterRoutes(c Context, r *chi.Mux) {
	uh := NewUserHandler(c)

	api := r.Route("/api", nil)

	api.Route("/users", func(r chi.Router) {
		r.Post("/", uh.registerHandlerFunc())
		r.Post("/login", uh.loginHandlerFunc())
	})

	api.Route("/user", func(r chi.Router) {
		r.Use(middleware.Verifier(c.jwt))
		r.Use(middleware.Authenticator)

		r.Get("/", uh.getHandlerFunc())

		r.Put("/", uh.profileHandlerFunc())
	})

	api.Route("/profiles", func(r chi.Router) {
		// public
		r.Get("/{username}", uh.profileHandlerFunc())

		// auth required
		auth := r.With(middleware.Verifier(c.jwt), middleware.Authenticator)
		auth.Post("/{username}/follow", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("not implemented yet")
		})

		auth.Delete("/{username}/follow", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("not implemented yet")
		})
	})

	api.Route("/articles", func(r chi.Router) {
		// public
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Articles not implemented yet")
		})

		r.Get("/{slug}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("GetArticle not implemented yet")
		})

		r.Get("/{slug}/comments", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("GetComments not implemented yet")
		})

		// auth required
		auth := r.With(middleware.Verifier(c.jwt), middleware.Authenticator)

		auth.Post("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("CreateArticle not implemented yet")
		})
		auth.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Feed not implemented yet")
		})
		auth.Put("/{slug}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("UpdateArticle not implemented yet")
		})
		auth.Delete("/{slug}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("DeleteArticle not implemented yet")
		})
		auth.Post("/{slug}/comments", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("AddComment not implemented yet")
		})
		auth.Delete("/{slug}/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("DeleteComment not implemented yet")
		})
		auth.Post("/{slug}/favorite", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Favorite not implemented yet")
		})
		auth.Delete("/{slug}/favorite", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Unfavorite not implemented yet")
		})

	})

	api.Get("/tags", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Tags not implemented yet")
	})
}
