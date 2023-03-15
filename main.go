package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rikstam/ent-blog-example/ent"
	"github.com/rikstam/ent-blog-example/ent/user"
	"html/template"
	"log"
	"net/http"
)

var (
	//go:embed templates/*
	resources embed.FS
	tmpl      = template.Must(template.ParseFS(resources, "templates/*"))
)

type server struct {
	client *ent.Client
}

func newServer(client *ent.Client) *server {
	return &server{client: client}
}

// newRouter creates a new router with the blog handlers mounted.
func newRouter(srv *server) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", srv.index)
	return r
} // index serves the blog home page

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	posts, err := s.client.Post.
		Query().
		WithAuthor().
		All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// read connection string from a CLI flag
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "database DSN")
	flag.Parse()

	// Instantiate the Ent client
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed connecting to Mysql: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	// seed DB if no posts
	if !client.Post.Query().ExistX(ctx) {
		if err := seed(ctx, client); err != nil {
			log.Fatalf("failed seeding the database: %v", err)
		}
	}

	srv := newServer(client)
	r := newRouter(srv)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func seed(ctx context.Context, client *ent.Client) error {
	// Check if user "Rikstam" exists
	r, err := client.User.Query().
		Where(
			user.Name("Rikstam"),
		).
		Only(ctx)

	switch {
	// If not, create the user.
	case ent.IsNotFound(err):
		r, err = client.User.Create().
			SetName("Rikstam").
			SetEmail("r@hello.world").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating user: %v", err)
		}
	case err != nil:
		return fmt.Errorf("failed querying user: %v", err)
	}
	// Finally, create a "Hello, world" blogpost.
	return client.Post.Create().
		SetTitle("Hello, World!").
		SetBody("This is my first post").
		SetAuthor(r).
		Exec(ctx)
}