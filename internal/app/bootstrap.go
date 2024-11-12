package app

import (
	"log"
	"os"

	"github.com/federicodosantos/socialize/internal/delivery/http"
	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/repository"
	"github.com/federicodosantos/socialize/internal/usecase"
	"github.com/federicodosantos/socialize/pkg/jwt"
	"github.com/federicodosantos/socialize/pkg/supabase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	supabaseStorage "github.com/supabase-community/storage-go"
)

type Bootstrap struct {
	db     *sqlx.DB
	router *chi.Mux
}

func NewBootstrap(db *sqlx.DB, router *chi.Mux) *Bootstrap {
	return &Bootstrap{
		db:     db,
		router: router,
	}
}

func (b *Bootstrap) InitApp() {
	// initialize jwt service
	jwtService, err := jwt.NewJwt(os.Getenv("JWT_SECRET_KEY"), os.Getenv("JWT_EXPIRED"))
	if err != nil {
		log.Printf("cannot initialize jwt service due to %s", err.Error())
	}

	// initialize supabase
	client := supabaseStorage.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"),
		map[string]string{
			"apikey": os.Getenv("SUPABASE_KEY"),
		})

	supabase := supabase.NewSupabaseStorage(client)

	// initialize middleware
	middleware := middleware.NewMiddleware(jwtService)

	// initialize repository
	userRepo := repository.NewUserRepo(b.db)

	// initialize usecase
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService, supabase)

	// init handler
	userHandler := http.NewUserHandler(userUsecase)

	b.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// init routes
	http.UserRoutes(b.router, userHandler, middleware)
}
