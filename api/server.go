package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/gu3sswho/simplebank/db/sqlc"
)

// Server serves HTTP requests for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer create server and setup routes
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.PUT("/accounts/:id", server.updateAccount)

	server.router = router

	return server
}

// Start run HTTP server on special address and port
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse is error wrapper
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
