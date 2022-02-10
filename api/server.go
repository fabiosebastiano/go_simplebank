package api

import (
	db "github.com/fabiosebastiano/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

//NewServer crea un nuovo server HTTP e prepara il routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	//router.PATCH("/accounts/:id", server.updateAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router

	return server
}

//Start il server su indirizzo specificato
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
