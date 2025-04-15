package api

import (
	"apps/ecosystem/core/handlers"
	"apps/ecosystem/core/stores"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer {
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"}, // does Set-Cookie even need to be here?
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	baseRouter := chi.NewRouter()
	router.Mount("/api/v1", baseRouter)
	
	tableStore := stores.NewTableStore(s.db)
	recordStore := stores.NewRecordStore(s.db)
	userStore := stores.NewUserStore(s.db)
	roleStore := stores.NewRoleStore(s.db)

	tableHandler := handlers.NewTableHandler(tableStore)
	recordHandler := handlers.NewRecordHandler(recordStore)
	userHandler := handlers.NewUserHandler(userStore, roleStore)

	tableHandler.RegisterRoutes(baseRouter)
	recordHandler.RegisterRoutes(baseRouter)
	userHandler.RegisterRoutes(baseRouter)

	log.Println("listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
