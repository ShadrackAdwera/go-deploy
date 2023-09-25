package api

import (
	"fmt"
	db "go-k8s/internal/db/sqlc"
	"go-k8s/internal/token"
	"go-k8s/internal/workers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	store       db.TxStore
	pasetoMaker token.Maker
	distro      workers.TaskDistributor
}

func NewServer(store db.TxStore, pasetoMaker token.Maker, distro workers.TaskDistributor) *Server {
	srv := Server{
		store:       store,
		pasetoMaker: pasetoMaker,
		distro:      distro,
	}

	router := gin.Default()

	// add cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"DELETE", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
	}))

	router.POST("/api/sign-up", srv.createTenantTx)
	router.POST("/api/login", srv.login)

	srv.router = router
	return &srv
}

func (srv *Server) StartServer(add string) error {
	return srv.router.Run(add)
}

func errJSON(err error) gin.H {
	return gin.H{"message": fmt.Errorf("error %w", err)}
}
