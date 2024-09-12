package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/transfers", server.createTransfer)
	router.POST("/accounts", server.createAccount)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
