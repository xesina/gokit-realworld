package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"net/http"
)

func RegisterRoutes(c Context, r *chi.Mux) {
	uh := NewUserHandler(c)
	ah := NewArticleHandler(c)

	api := r.Route("/api", nil)

	// Always parse token if available
	api.Use(middleware.Verifier(c.jwt))

	api.Route("/users", func(r chi.Router) {
		r.Post("/", uh.registerHandlerFunc())
		r.Post("/login", uh.loginHandlerFunc())
	})

	api.Route("/user", func(r chi.Router) {
		r.Use(middleware.Authenticator)
		r.Get("/", uh.getHandlerFunc())
		r.Put("/", uh.updateHandlerFunc())
	})

	api.Route("/profiles", func(r chi.Router) {
		// public
		r.Get("/{username}", uh.profileHandlerFunc())

		// auth required
		auth := r.With(middleware.Authenticator)
		auth.Post("/{username}/follow", uh.followHandlerFunc())
		auth.Delete("/{username}/follow", uh.unfollowHandlerFunc())
	})

	api.Route("/articles", func(r chi.Router) {
		// public
		r.Get("/", ah.listHandlerFunc())

		r.Get("/{slug}", ah.getHandlerFunc())

		r.Get("/{slug}/comments", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("GetComments not implemented yet")
		})

		// auth required
		auth := r.With(middleware.Authenticator)

		auth.Post("/", ah.createHandlerFunc())
		auth.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Feed not implemented yet")
		})
		auth.Put("/{slug}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("UpdateArticle not implemented yet")
		})
		auth.Delete("/{slug}", ah.deleteHandlerFunc())
		auth.Post("/{slug}/comments", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("AddComment not implemented yet")
		})
		auth.Delete("/{slug}/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("DeleteComment not implemented yet")
		})
		auth.Post("/{slug}/favorite", ah.favoriteHandlerFunc())
		auth.Delete("/{slug}/favorite", ah.unfavoriteHandlerFunc())

	})

	api.Get("/tags", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Tags not implemented yet")
	})
}
