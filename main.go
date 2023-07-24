package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// App contains the state of the application
type App struct {
	dataDir string
	db      *sqlx.DB
}

func (s App) Close() {
	s.db.Close()
}

var app App

func main() {
	// app config
	var (
		dataDir = "data"
		dsn     = "root:@/imgram"
	)

	if val, ok := os.LookupEnv("IMGRAM_DATA_DIR"); ok {
		dataDir = val
	}
	if val, ok := os.LookupEnv("IMGRAM_DSN"); ok {
		dsn = val
	}

	// init database
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Printf("database: open error: %v\n", err)
		panic(err)
	}

	// init app state
	app = App{
		dataDir: dataDir,
		db:      db,
	}
	defer app.Close()

	// init handlers
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Healthcheck
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/posts", func(r chi.Router) {
		r.Use(UserAuth)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// list posts with last 2 comments for each post
		})
		r.Post("/", PostImage)
		r.Post("/{post_id}/comments", CommentPost)
		r.Delete("/{post_id}/comments/{comment_id}", DeleteComment)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v\n", err)
		}
	}
}
