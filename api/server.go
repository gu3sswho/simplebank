package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/gu3sswho/simplebank/db/sqlc"
	"github.com/gu3sswho/simplebank/token"
	"github.com/gu3sswho/simplebank/util"
)

// Server serves HTTP requests for banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer create server and setup routes
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// setupRouter sets all routers for the server
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.GET("/users/:username", server.getUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.PUT("/accounts/:id", server.updateAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start run HTTP server on special address and port
func (server *Server) Start() error {
	return server.router.Run(server.config.ServerAddr)
}

// errorResponse is error wrapper
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
