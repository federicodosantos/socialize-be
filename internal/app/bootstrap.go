package app

import (
	"log"
	"os"

	httpHandler "github.com/federicodosantos/socialize/internal/delivery/http"
	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/repository"
	"github.com/federicodosantos/socialize/internal/usecase"
	"github.com/federicodosantos/socialize/pkg/jwt"
	"github.com/federicodosantos/socialize/pkg/supabase"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	supabaseStorage "github.com/supabase-community/storage-go"
)

type Bootstrap struct {
	db     *sqlx.DB
	router *chi.Mux
	logger *zap.SugaredLogger
}

func NewBootstrap(db *sqlx.DB, router *chi.Mux, logger *zap.SugaredLogger) *Bootstrap {
	return &Bootstrap{
		db:     db,
		router: router,
		logger: logger,
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

	// initialize repository
	userRepo := repository.NewUserRepo(b.db)
	postRepo := repository.NewPostRepo(b.db)
	commentRepo := repository.NewCommentRepo(b.db)

	// initialize usecase
	fileUsecase := usecase.NewFileUsecase(supabase)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)
	postUsecase := usecase.NewPostUsecase(postRepo, commentRepo)

	// init handler
	fileHandler := httpHandler.NewFileHandler(fileUsecase)
	userHandler := httpHandler.NewUserHandler(userUsecase)
	postHandler := httpHandler.NewPostHandler(postUsecase)

	// initialize middleware
	middleware := middleware.NewMiddleware(jwtService, b.logger)

	b.router.Use(middleware.LoggingMiddleware)

	b.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// init routes
	httpHandler.FileRoutes(b.router, fileHandler, middleware)
	httpHandler.UserRoutes(b.router, userHandler, middleware)
	httpHandler.PostRoutes(b.router, postHandler, middleware)

	//health check
	util.HealthCheck(b.router, b.db)
}
