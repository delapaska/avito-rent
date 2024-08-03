package api

import (
	"database/sql"

	"github.com/delapaska/avito-rent/configs"
	"github.com/delapaska/avito-rent/service/auth"
	dummyauth "github.com/delapaska/avito-rent/service/dummyAuth"
	"github.com/delapaska/avito-rent/service/flat"
	"github.com/delapaska/avito-rent/service/house"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	addr   string
	engine *gin.Engine
}

func NewAPIServer(db *sql.DB) *APIServer {

	engine := gin.New()
	dummyStore := dummyauth.NewStore(db)
	dummyHandler := dummyauth.NewHandler(dummyStore)
	dummyHandler.RegisterRoutes(engine)

	houseStore := house.NewStore(db)
	houseHandler := house.NewHandler(houseStore)
	houseHandler.RegisterRoutes(engine)

	flatStore := flat.NewStore(db)
	flatHandler := flat.NewHandler(flatStore)
	flatHandler.RegisterRoutes(engine)

	authStore := auth.NewStore(db)
	authHandler := auth.NewHandler(authStore)
	authHandler.RegisterRoutes(engine)

	return &APIServer{
		addr:   ":" + configs.Envs.Port,
		engine: engine,
	}
}

func (s *APIServer) Run() {

	s.engine.Run(s.addr)
}
